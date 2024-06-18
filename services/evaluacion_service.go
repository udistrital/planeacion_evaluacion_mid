package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	evaluacionhelper "github.com/udistrital/planeacion_evaluacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
	"golang.org/x/sync/errgroup"
)

const (
	CodigoEstadoPlan                     string = "A_SP" //6153355601c7a2365b2fb2a1
	CodigoEstadoSeguimiento              string = "AV"   //622ba49216511e93a95c326d
	CodigoTipoSeguimiento                string = "S_SP" //61f236f525e40c582a0840d0
	ABREVIACION_AVALADO_PARA_SEGUIMIENTO string = "AV"
	ABREVIACION_SEGUIMIENTO_PLAN_ACCION  string = "S_SP"
)

var NOMBRE_TRIMESTRE = map[int]string{
	1: "Trimestre Uno",
	2: "Trimestre Dos",
	3: "Trimestre Tres",
	4: "Trimestre Cuatro",
}

func GetPlanesPeriodo(vigencia string, unidad string) (respuesta []map[string]interface{}, outputError map[string]interface{}) {

	if len(vigencia) == 0 || len(unidad) == 0 {
		outputError = map[string]interface{}{
			"funcion": "error al decodificar el cuerpo de la solicitud",
			"status":  "400",
		}
		return nil, outputError
	}

	var resPlan map[string]interface{}
	var resSeguimiento map[string]interface{}
	respuesta = make([]map[string]interface{}, 0)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan?query=estado_plan_id:6153355601c7a2365b2fb2a1,dependencia_id:`+unidad+`,vigencia:`+vigencia, &resPlan); err == nil {
		planes := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(resPlan, &planes)
		// formatdata.JsonPrint(planes)
		if fmt.Sprintf("%v", planes) == "[]" {
			outputError = map[string]interface{}{
				"funcion": "No se tienen planes en seguimiento para la dependencia y la vigencia",
				"err":     err,
				"status":  "400",
			}
		}

		periodos := GetPeriodos(vigencia, true)
		trimestres := GetTrimestres(vigencia)

		for _, plan := range planes {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,estado_seguimiento_id:622ba49216511e93a95c326d,plan_id:`+plan["_id"].(string), &resSeguimiento); err == nil {
				seguimientos := make([]map[string]interface{}, 1)
				request.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)
				if fmt.Sprintf("%v", seguimientos) == "[]" {
					continue
				}

				var periodosSelecionados []map[string]interface{}
				for _, seguimiento := range seguimientos {
					for _, periodo := range periodos {
						if seguimiento["periodo_seguimiento_id"] == periodo["_id"] {
							for _, trimestre := range trimestres {
								var trimestreId float64
								if reflect.TypeOf(trimestre["Id"]).String() == "string" {
									trimestreId, _ = strconv.ParseFloat(trimestre["Id"].(string), 64)
								} else {
									trimestreId = trimestre["Id"].(float64)
								}
								var periodoId float64
								if reflect.TypeOf(periodo["periodo_id"]).String() == "string" {
									periodoId, _ = strconv.ParseFloat(periodo["periodo_id"].(string), 64)
								} else {
									periodoId = periodo["periodo_id"].(float64)
								}

								if trimestreId == periodoId {
									periodosSelecionados = append(periodosSelecionados, map[string]interface{}{"nombre": trimestre["ParametroId"].(map[string]interface{})["Nombre"].(string), "id": periodo["_id"]})
									break
								}
							}
							break
						}
					}
				}
				respuesta = append(respuesta, map[string]interface{}{"plan": plan["nombre"], "id": plan["_id"], "periodos": periodosSelecionados})
			} else {
				outputError = map[string]interface{}{
					"err":    err,
					"status": "400",
				}
			}
		}
	} else {
		outputError = map[string]interface{}{
			"err":    err,
			"status": "400",
		}
	}
	return respuesta, outputError
}

func GetEvaluacion(vigencia string, plan string, periodoId string) (interface{}, error) {

	var evaluacion []map[string]interface{}

	if len(vigencia) == 0 || len(plan) == 0 || len(periodoId) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	if len(plan) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	if len(periodoId) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}

	trimestres := GetPeriodos(vigencia, true)

	if len(trimestres) < 4 {
		return nil, nil
	} else {
		i := 0
		for index, periodo := range trimestres {
			if periodo["_id"] == periodoId {
				i = index
				break
			}
		}

		var periodo_id interface{}

		for _, periodo := range trimestres {
			if periodo["_id"] == periodoId {
				periodo_id = periodo["periodo_id"]
				break
			}
		}
		evaluacion = GetEvaluacionInterno(plan, trimestres, i, vigencia, periodo_id)
		return evaluacion, nil
	}
}

func Unidades(plan string, vigencia string) (interface{}, error) {
	if nombrePlan, err := url.QueryUnescape(plan); err == nil {
		if data, err := GetUnidadesPorPlanYVigencia(nombrePlan, vigencia); err == nil {
			return data, nil
		} else {
			return nil, errors.New("Error obteniendo las unidades del plan y la vigencia dados ")
		}
	} else {
		return nil, errors.New("Error obteniendo las unidades del plan y la vigencia dados " + err.Error())

	}
}

func GetUnidadesPorPlanYVigencia(nombrePlan string, vigencia string) (unidades []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetUnidadesPorPlanYVigencia",
				"err":     err,
				"status":  "400",
			}
		}
	}()

	var respuestaEstado map[string]interface{}
	var respuestaTipoSeguimiento map[string]interface{}
	var respuestaSeguimiento map[string]interface{}
	var respuestaTipoDependencia []map[string]interface{}

	var estadoSeguimiento []map[string]interface{}
	var tipoSeguimiento []map[string]interface{}
	idsUnidades := make([]string, 0)
	unidades = make([]map[string]interface{}, 0)
	var mutex sync.Mutex
	wge := new(errgroup.Group)
	wge2 := new(errgroup.Group)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_AVALADO_PARA_SEGUIMIENTO, &respuestaEstado); err != nil {
		outputError = map[string]interface{}{
			"err":    err,
			"status": "404",
		}
	}
	request.LimpiezaRespuestaRefactor(respuestaEstado, &estadoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_SEGUIMIENTO_PLAN_ACCION, &respuestaTipoSeguimiento); err != nil {
		outputError = map[string]interface{}{
			"err":    err,
			"status": "404",
		}
	}
	request.LimpiezaRespuestaRefactor(respuestaTipoSeguimiento, &tipoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:`+tipoSeguimiento[0]["_id"].(string)+`,estado_seguimiento_id:`+estadoSeguimiento[0]["_id"].(string), &respuestaSeguimiento); err == nil {
		var seguimientos []map[string]interface{}
		request.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimientos)
		for _, seguimiento := range seguimientos {
			seguimiento := seguimiento
			wge.Go(func() error {
				// Esta en los planes que ya se trajeron?
				var respuestaPlan map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan/`+seguimiento["plan_id"].(string), &respuestaPlan); err == nil {
					var plan map[string]interface{}
					request.LimpiezaRespuestaRefactor(respuestaPlan, &plan)

					if plan["nombre"] == nombrePlan && vigencia == plan["vigencia"] {
						existeIdUnidad := false
						for _, idUnidad := range idsUnidades {
							if idUnidad == plan["dependencia_id"].(string) {
								existeIdUnidad = true
							}
						}
						if !existeIdUnidad {
							mutex.Lock()
							idsUnidades = append(idsUnidades, plan["dependencia_id"].(string))
							mutex.Unlock()
						}
					}

				} else {
					outputError = map[string]interface{}{
						"err":    err,
						"status": "404",
					}
				}
				return nil
			})
		}
		if err := wge.Wait(); err != nil {
			return nil, outputError
		}
		if len(idsUnidades) > 0 {
			valor := 0
			for _, idUnidad := range idsUnidades {
				idUnidad := idUnidad
				wge2.Go(func() error {
					valor++
					if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId__Id:"+idUnidad, &respuestaTipoDependencia); err == nil {
						aux := respuestaTipoDependencia[0]["DependenciaId"].(map[string]interface{})
						delete(aux, "DependenciaTipoDependencia")
						aux["TipoDependencia"] = respuestaTipoDependencia[0]["TipoDependenciaId"]
						mutex.Lock()
						unidades = append(unidades, aux)
						mutex.Unlock()
						respuestaTipoDependencia = nil
					} else {
						outputError = map[string]interface{}{
							"err":    err,
							"status": "404",
						}
					}
					return nil
				})
			}
			if err := wge2.Wait(); err != nil {
				return nil, outputError
			}
		}
	} else {
		outputError = map[string]interface{}{
			"err":    err,
			"status": "404",
		}
	}
	return unidades, outputError
}

func Avances(plan string, vigencia string, unidad string) (interface{}, error) {
	if nombrePlan, err1 := url.QueryUnescape(plan); err1 == nil {
		if data, err2 := GetAvances(nombrePlan, vigencia, unidad); err2 == nil {
			return data, nil
		} else {
			return nil, errors.New("Error obteniendo los avances 1 ")
		}
	} else {
		return nil, errors.New("Error obteniendo los avances 2" + err1.Error())
	}
}

func GetAvances(nombrePlan string, idVigencia string, idUnidad string) (respuesta map[string]interface{}, outputError map[string]interface{}) {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		outputError = map[string]interface{}{
	// 			"funcion": "GetAvances",
	// 			"err":     err,
	// 			"status":  "404",
	// 		}
	// 		panic(outputError)
	// 	}
	// }()
	respuesta = make(map[string]interface{}, 0)
	avance := map[int]float64{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
	}

	if planes, err := GetPlanesPeriodo(idVigencia, idUnidad); err != nil {
		outputError = map[string]interface{}{
			"funcion": "GetAvances",
			"err":     err,
			"status":  "404",
		}
		return nil, outputError
	} else {
		planes = evaluacionhelper.FiltrarArreglo(planes, func(plan map[string]interface{}) bool {
			return plan["plan"] == nombrePlan
		})
		plan := planes[0]
		respuesta["plan"] = map[string]interface{}{
			"id":     plan["id"],
			"nombre": plan["plan"],
		}
		periodos := plan["periodos"].([]map[string]interface{})
		ultimoPeriodo := periodos[len(periodos)-1]
		respuesta["periodo"] = map[string]interface{}{
			"id":     ultimoPeriodo["id"],
			"nombre": ultimoPeriodo["nombre"],
		}

		trimestres := GetPeriodos(idVigencia, true)
		var periodo_id interface{}

		// formatdata.JsonPrint(trimestres)
		if len(trimestres) != 0 {
			for index, trimestre := range trimestres {
				for _, periodo := range periodos {
					if trimestre["_id"] == periodo["id"] {
						if trimestre["_id"] == periodo["id"] {
							periodo_id = trimestre["periodo_id"]
						}
						actividades := GetEvaluacionInterno(plan["id"].(string), trimestres, index, idVigencia, periodo_id)
						var numeroActividad = "0"
						for _, actividad := range actividades {
							if numeroActividad != actividad["numero"] {
								// Se guarda el número de la actividad, si hay más registros con la misma actividad no se sumaran
								numeroActividad = actividad["numero"].(string)
								for i := 1; i < 5; i++ {
									infoTrimestre := actividad["trimestre"+strconv.Itoa(i)].(map[string]interface{})
									if fmt.Sprintf("%v", infoTrimestre) != "map[]" && periodo["nombre"] == NOMBRE_TRIMESTRE[i] {
										// Se tiene info del Trimestre y se está evaluando el mismo periodo
										avanceActividad := infoTrimestre["actividad"].(float64)
										if avanceActividad >= 1 {
											avance[i] = avance[i] + actividad["ponderado"].(float64)
										} else {
											avance[i] = avance[i] + actividad["ponderado"].(float64)*float64(avanceActividad)
										}
									}
								}
							}
						}

					}
				}
			}
		}

	}
	respuesta["Trimestres"] = avance
	respuesta["Promedio"] = (avance[1] + avance[2] + avance[3] + avance[4]) / 4
	return respuesta, outputError
}

func GetPlanesPeriodoInterno(unidad string, vigencia string) (respuesta []map[string]interface{}, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}
	}()

	var resPlan map[string]interface{}
	var resSeguimiento map[string]interface{}
	respuesta = make([]map[string]interface{}, 0)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan?query=estado_plan_id:6153355601c7a2365b2fb2a1,dependencia_id:`+unidad+`,vigencia:`+vigencia, &resPlan); err == nil {
		planes := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(resPlan, &planes)
		// formatdata.JsonPrint(planes)
		if fmt.Sprintf("%v", planes) == "[]" {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}

		periodos := GetPeriodos(vigencia, true)

		trimestres := GetTrimestres(vigencia)

		for _, plan := range planes {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:61f236f525e40c582a0840d0,estado_seguimiento_id:622ba49216511e93a95c326d,plan_id:`+plan["_id"].(string), &resSeguimiento); err == nil {
				seguimientos := make([]map[string]interface{}, 1)
				request.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)
				if fmt.Sprintf("%v", seguimientos) == "[]" {
					continue
				}

				var periodosSelecionados []map[string]interface{}
				for _, seguimiento := range seguimientos {
					for _, periodo := range periodos {
						if seguimiento["periodo_seguimiento_id"] == periodo["_id"] {
							for _, trimestre := range trimestres {
								var trimestreId float64
								if reflect.TypeOf(trimestre["Id"]).String() == "string" {
									trimestreId, _ = strconv.ParseFloat(trimestre["Id"].(string), 64)
								} else {
									trimestreId = trimestre["Id"].(float64)
								}
								var periodoId float64
								if reflect.TypeOf(periodo["periodo_id"]).String() == "string" {
									periodoId, _ = strconv.ParseFloat(periodo["periodo_id"].(string), 64)
								} else {
									periodoId = periodo["periodo_id"].(float64)
								}

								if trimestreId == periodoId {
									periodosSelecionados = append(periodosSelecionados, map[string]interface{}{"nombre": trimestre["ParametroId"].(map[string]interface{})["Nombre"].(string), "id": periodo["_id"]})
									break
								}
							}
							break
						}
					}
				}
				respuesta = append(respuesta, map[string]interface{}{"plan": plan["nombre"], "id": plan["_id"], "periodos": periodosSelecionados})
			} else {
				outputError = errors.New("error al decodificar el cuerpo de la solicitud")
			}
		}
	} else {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
	}
	return respuesta, outputError
}

func PlanesAEvaluar() (planes []string, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}
	}()

	var respuestaEstado map[string]interface{}
	var respuestaTipoSeguimiento map[string]interface{}
	var respuestaSeguimiento map[string]interface{}
	var estadoSeguimiento []map[string]interface{}
	var tipoSeguimiento []map[string]interface{}
	var mutex sync.Mutex
	wge := new(errgroup.Group)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_AVALADO_PARA_SEGUIMIENTO, &respuestaEstado); err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	request.LimpiezaRespuestaRefactor(respuestaEstado, &estadoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_SEGUIMIENTO_PLAN_ACCION, &respuestaTipoSeguimiento); err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	request.LimpiezaRespuestaRefactor(respuestaTipoSeguimiento, &tipoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:`+tipoSeguimiento[0]["_id"].(string)+`,estado_seguimiento_id:`+estadoSeguimiento[0]["_id"].(string)+`,activo:true`, &respuestaSeguimiento); err == nil {
		var seguimientos []map[string]interface{}
		request.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimientos)
		for _, seguimiento := range seguimientos {
			seguimiento := seguimiento
			wge.Go(func() error {
				// Esta en los planes que ya se trajeron?
				var respuestaPlan map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan/`+seguimiento["plan_id"].(string), &respuestaPlan); err == nil {
					var plan map[string]interface{}
					existeNombrePlan := false
					request.LimpiezaRespuestaRefactor(respuestaPlan, &plan)
					for _, nombre := range planes {
						if nombre == plan["nombre"].(string) {
							existeNombrePlan = true
						}
					}
					if !existeNombrePlan {
						mutex.Lock()
						planes = append(planes, plan["nombre"].(string))
						mutex.Unlock()
					}
				} else {
					outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
				}
				return nil
			})
		}
		if err := wge.Wait(); err != nil {
			return nil, errors.New(err.Error())
		}
	} else {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}

	return planes, outputError
}

func GetEvaluacionInterno(planId string, periodos []map[string]interface{}, trimestre int, vigencia string, periodo_id interface{}) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	var trimestrenombreL []map[string]interface{}
	var trimestrenombre map[string]interface{}

	detalle := make(map[string]interface{})
	actividades := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=estado_seguimiento_id:622ba49216511e93a95c326d,plan_id:`+planId+`,periodo_seguimiento_id:`+periodos[trimestre]["_id"].(string), &resSeguimiento); err == nil {
		aux := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(resSeguimiento, &aux)
		if fmt.Sprintf("%v", aux) == "[]" {
			return nil
		}
		seguimiento = aux[0]
		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &actividades)
		for actividadId, act := range actividades {
			id, segregado := actividades[actividadId].(map[string]interface{})["id"].(string)
			var actividad map[string]interface{}

			if segregado && id != "" {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id, &resSeguimientoDetalle); err == nil {
					request.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
					actividad = evaluacionhelper.ConvertirStringJson(detalle)
				}
			} else {
				actividad = act.(map[string]interface{})
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",Id:"+periodo_id.(string)+"", &trimestrenombre); err == nil {
				request.LimpiezaRespuestaRefactor(trimestrenombre, &trimestrenombreL)
			} else {
				err = errors.New("error al decodificar el cuerpo de la solicitud: 404" + err.Error())
			}

			// Variable para guardar el nombre del trimestre
			var NombreTrim string

			// Acceder al campo "Nombre" en el mapa "ParametroId"
			for _, item := range trimestrenombreL {
				if param, ok := item["ParametroId"].(map[string]interface{}); ok {
					if nombre, ok := param["Nombre"].(string); ok {
						NombreTrim = nombre
						break
					}
				}
			}

			for indexPeriodo, periodo := range periodos {
				if indexPeriodo > trimestre {
					break
				}
				resIndicadores := GetEvaluacionTrimestre(planId, periodo["_id"].(string), actividadId)

				for _, resIndicador := range resIndicadores {
					indice := -1
					for index, eval := range evaluacion {
						if eval["numero"] == actividad["informacion"].(map[string]interface{})["index"] && eval["indicador"] == resIndicador["indicador"] {
							indice = index
							break
						}

					}

					var trimestreNom string
					if NombreTrim == "Trimestre Uno" {
						trimestreNom = "trimestre1"
					} else if NombreTrim == "Trimestre Dos" {
						trimestreNom = "trimestre2"
					} else if NombreTrim == "Trimestre Tres" {
						trimestreNom = "trimestre3"
					} else if NombreTrim == "Trimestre Cuatro" {
						trimestreNom = "trimestre4"
					}

					if indice == -1 {
						evaluacionAct := map[string]interface{}{
							"actividad":  actividad["informacion"].(map[string]interface{})["descripcion"],
							"numero":     actividad["informacion"].(map[string]interface{})["index"],
							"periodo":    actividad["informacion"].(map[string]interface{})["periodo"],
							"ponderado":  actividad["informacion"].(map[string]interface{})["ponderacion"],
							"trimestre1": make(map[string]interface{}),
							"trimestre2": make(map[string]interface{}),
							"trimestre3": make(map[string]interface{}),
							"trimestre4": make(map[string]interface{}),
						}
						evaluacionAct["indicador"] = resIndicador["indicador"]
						evaluacionAct["unidad"] = resIndicador["unidad"]
						evaluacionAct["formula"] = resIndicador["formula"]
						evaluacionAct["meta"] = resIndicador["metaA"].(float64)
						evaluacionAct[trimestreNom] = map[string]interface{}{
							"acumulado":            resIndicador["acumulado"],
							"denominador":          resIndicador["denominador"],
							"meta":                 resIndicador["meta"],
							"numerador":            resIndicador["numerador"],
							"periodo":              resIndicador["periodo"],
							"numeradorAcumulado":   resIndicador["numeradorAcumulado"],
							"denominadorAcumulado": resIndicador["denominadorAcumulado"],
							"brecha":               resIndicador["brecha"],
						}

						evaluacion = append(evaluacion, evaluacionAct)
					} else {
						evaluacion[indice][trimestreNom] = map[string]interface{}{
							"acumulado":            resIndicador["acumulado"],
							"denominador":          resIndicador["denominador"],
							"meta":                 resIndicador["meta"],
							"numerador":            resIndicador["numerador"],
							"periodo":              resIndicador["periodo"],
							"numeradorAcumulado":   resIndicador["numeradorAcumulado"],
							"denominadorAcumulado": resIndicador["denominadorAcumulado"],
							"brecha":               resIndicador["brecha"],
						}
					}
				}
			}
		}

		request.SortSlice(&evaluacion, "numero")
		agrupacion_actividades := make(map[string][]int)
		for i, eval := range evaluacion {
			if _, ok := agrupacion_actividades[eval["numero"].(string)]; !ok {
				agrupacion_actividades[eval["numero"].(string)] = []int{}
			}
			agrupacion_actividades[eval["numero"].(string)] = append(agrupacion_actividades[eval["numero"].(string)], i)
		}

		for _, idxs := range agrupacion_actividades {
			sum1 := 0.0
			sum2 := 0.0
			sum3 := 0.0
			sum4 := 0.0

			for _, i := range idxs {
				if fmt.Sprintf("%v", evaluacion[i]["trimestre1"]) != "map[]" {
					if evaluacion[i]["trimestre1"].(map[string]interface{})["meta"].(float64) > 1 {
						sum1 = sum1 + 1.0
					} else {
						sum1 = sum1 + evaluacion[i]["trimestre1"].(map[string]interface{})["meta"].(float64)
					}
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre2"]) != "map[]" {
					if evaluacion[i]["trimestre2"].(map[string]interface{})["meta"].(float64) > 1 {
						sum2 = sum2 + 1.0
					} else {
						sum2 = sum2 + evaluacion[i]["trimestre2"].(map[string]interface{})["meta"].(float64)
					}
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre3"]) != "map[]" {
					if evaluacion[i]["trimestre3"].(map[string]interface{})["meta"].(float64) > 1 {
						sum3 = sum3 + 1.0
					} else {
						sum3 = sum3 + evaluacion[i]["trimestre3"].(map[string]interface{})["meta"].(float64)
					}
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre4"]) != "map[]" {
					if evaluacion[i]["trimestre4"].(map[string]interface{})["meta"].(float64) > 1 {
						sum4 = sum4 + 1.0
					} else {
						sum4 = sum4 + evaluacion[i]["trimestre4"].(map[string]interface{})["meta"].(float64)
					}
				}
			}

			cont := len(idxs)
			cumplActividad1 := math.Floor((sum1/float64(cont))*1000) / 1000
			cumplActividad2 := math.Floor((sum2/float64(cont))*1000) / 1000
			cumplActividad3 := math.Floor((sum3/float64(cont))*1000) / 1000
			cumplActividad4 := math.Floor((sum4/float64(cont))*1000) / 1000

			if cumplActividad1 > 1 {
				cumplActividad1 = 1
			}

			if cumplActividad2 > 1 {
				cumplActividad2 = 1
			}

			if cumplActividad3 > 1 {
				cumplActividad3 = 1
			}

			if cumplActividad4 > 1 {
				cumplActividad4 = 1
			}

			for _, i := range idxs {
				if fmt.Sprintf("%v", evaluacion[i]["trimestre1"]) != "map[]" {
					evaluacion[i]["trimestre1"].(map[string]interface{})["actividad"] = cumplActividad1
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre2"]) != "map[]" {
					evaluacion[i]["trimestre2"].(map[string]interface{})["actividad"] = cumplActividad2
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre3"]) != "map[]" {
					evaluacion[i]["trimestre3"].(map[string]interface{})["actividad"] = cumplActividad3
				}
				if fmt.Sprintf("%v", evaluacion[i]["trimestre4"]) != "map[]" {
					evaluacion[i]["trimestre4"].(map[string]interface{})["actividad"] = cumplActividad4
				}
			}
		}

		return evaluacion
	}
	return nil
}

func GetEvaluacionTrimestre(planId string, periodoId string, actividadId string) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	actividades := make(map[string]interface{})
	detalle := make(map[string]interface{})

	idEstadoSeguimiento, err := evaluacionhelper.GetIdCodigoAbreviacion("estado-seguimiento", CodigoEstadoSeguimiento)
	if err != nil {
		return nil
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=estado_seguimiento_id:`+idEstadoSeguimiento+`,plan_id:`+planId+`,periodo_seguimiento_id:`+periodoId, &resSeguimiento); err == nil {
		aux := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(resSeguimiento, &aux)
		if fmt.Sprintf("%v", aux) == "[]" {
			return nil
		}

		seguimiento = aux[0]

		datoStr := seguimiento["dato"].(string)
		json.Unmarshal([]byte(datoStr), &actividades)

		if actividades[actividadId] == nil {
			return nil
		}

		var indicadores []interface{}
		var resultados []interface{}
		id, segregado := actividades[actividadId].(map[string]interface{})["id"]

		if segregado && id != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id.(string), &resSeguimientoDetalle); err == nil {
				request.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
				detalle = evaluacionhelper.ConvertirStringJson(detalle)
				if fmt.Sprintf("%v", detalle["cuantitativo"]) != "map[]" {
					indicadores = detalle["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})
					resultados = detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})
				} else {
					indicadores = []interface{}{}
					resultados = []interface{}{}
				}
			}
		} else {
			indicadores = actividades[actividadId].(map[string]interface{})["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})
			resultados = actividades[actividadId].(map[string]interface{})["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})
		}

		for i := 0; i < len(indicadores); i++ {
			var metaA float64
			if indicadores[i].(map[string]interface{})["meta"] == nil {
				metaA = 0
			} else {
				if reflect.TypeOf(indicadores[i].(map[string]interface{})["meta"]).String() == "string" {
					metaA, _ = strconv.ParseFloat(indicadores[i].(map[string]interface{})["meta"].(string), 64)
				} else {
					metaA = indicadores[i].(map[string]interface{})["meta"].(float64)
				}
			}

			evaluacion = append(evaluacion, map[string]interface{}{
				"indicador":            indicadores[i].(map[string]interface{})["nombre"],
				"formula":              indicadores[i].(map[string]interface{})["formula"],
				"metaA":                metaA,
				"unidad":               indicadores[i].(map[string]interface{})["unidad"],
				"numerador":            indicadores[i].(map[string]interface{})["reporteNumerador"],
				"denominador":          indicadores[i].(map[string]interface{})["reporteDenominador"],
				"periodo":              resultados[i].(map[string]interface{})["indicador"],
				"acumulado":            resultados[i].(map[string]interface{})["indicadorAcumulado"],
				"meta":                 resultados[i].(map[string]interface{})["avanceAcumulado"],
				"numeradorAcumulado":   resultados[i].(map[string]interface{})["acumuladoNumerador"],
				"denominadorAcumulado": resultados[i].(map[string]interface{})["acumuladoDenominador"],
				"brecha":               resultados[i].(map[string]interface{})["brechaExistente"],
				"actividad":            0})
		}
		return evaluacion
	}
	return nil
}

func EvaluacionDetalle(vigencia string, plan string, periodoId string) (evaluacion []map[string]interface{}, outputError error) {

	defer func() {
		if err := recover(); err != nil {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}
	}()

	trimestres := GetPeriodos(vigencia, true)

	if len(trimestres) == 0 {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
	} else {
		i := 0
		for index, periodo := range trimestres {
			if periodo["_id"] == periodoId {
				i = index
				break
			}
		}

		evaluacion = GetEvaluacionInterno(plan, trimestres, i, vigencia, periodoId)

	}
	return evaluacion, outputError
}
func GetTrimestres(vigencia string) []map[string]interface{} {

	var res map[string]interface{}
	var trimestre []map[string]interface{}
	var trimestres []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T1", &res); err == nil {
		request.LimpiezaRespuestaRefactor(res, &trimestre)
		trimestres = append(trimestres, trimestre...)

		trimestre = nil
		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T2", &res); err == nil {
			request.LimpiezaRespuestaRefactor(res, &trimestre)
			trimestres = append(trimestres, trimestre...)

			trimestre = nil
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T3", &res); err == nil {
				request.LimpiezaRespuestaRefactor(res, &trimestre)
				trimestres = append(trimestres, trimestre...)

				trimestre = nil
				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T4", &res); err == nil {
					request.LimpiezaRespuestaRefactor(res, &trimestre)
					trimestres = append(trimestres, trimestre...)
				} else {
					err = errors.New("error al decodificar el cuerpo de la solicitud: GetTrimestres:400 " + err.Error())
				}
			} else {
				err = errors.New("error al decodificar el cuerpo de la solicitud: GetTrimestres:400 " + err.Error())
			}
		} else {
			err = errors.New("error al decodificar el cuerpo de la solicitud: GetTrimestres:400 " + err.Error())
		}
	} else {
		err = errors.New("error al decodificar el cuerpo de la solicitud: GetTrimestres:400 " + err.Error())
	}

	return trimestres
}

func GetPeriodos(vigencia string, modo bool) []map[string]interface{} {
	var periodos []map[string]interface{}
	var resPeriodo map[string]interface{}
	var wg sync.WaitGroup

	trimestres := GetTrimestres(vigencia)

	periodosMutex := sync.Mutex{}
	for _, trimestre := range trimestres {
		wg.Add(1)
		if fmt.Sprintf("%v", trimestre) == "map[]" {
			wg.Done()
			continue
		}
		go func(trimestreId int, wg *sync.WaitGroup, periodos *[]map[string]interface{}) {
			periodosMutex.Lock()
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?fields=_id,periodo_id&query=tipo_seguimiento_id:61f236f525e40c582a0840d0,periodo_id:`+strconv.Itoa(trimestreId), &resPeriodo); err == nil {
				var periodo []map[string]interface{}
				request.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
				if modo {
					(*periodos) = append((*periodos), periodo...)
				} else {
					(*periodos) = append((*periodos), periodo[len(periodo)-1])
				}

			}
			periodosMutex.Unlock()
			wg.Done()
		}(int(trimestre["Id"].(float64)), &wg, &periodos)
	}

	wg.Wait()

	request.SortSlice(&periodos, "periodo_id")
	return periodos
}

func PlanDetalle(vigencia string, unidad string) (result []map[string]interface{}, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}
	}()

	var resPlan map[string]interface{}
	var resSeguimiento map[string]interface{}

	idEstadoPlan, err := evaluacionhelper.GetIdCodigoAbreviacion("estado-plan", CodigoEstadoPlan)
	if err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan?query=estado_plan_id:`+idEstadoPlan+`,dependencia_id:`+unidad+`,vigencia:`+vigencia, &resPlan); err == nil {
		planes := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(resPlan, &planes)
		if fmt.Sprintf("%v", planes) == "[]" {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}

		periodos := GetPeriodos(vigencia, true)

		trimestres := GetTrimestres(vigencia)

		idEstadoSeguimiento, err1 := evaluacionhelper.GetIdCodigoAbreviacion("estado-seguimiento", CodigoEstadoSeguimiento)
		idTipoSeguimiento, err2 := evaluacionhelper.GetIdCodigoAbreviacion("tipo-seguimiento", CodigoTipoSeguimiento)
		if err1 != nil || err2 != nil {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}

		for _, plan := range planes {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:`+idTipoSeguimiento+`,estado_seguimiento_id:`+idEstadoSeguimiento+`,plan_id:`+plan["_id"].(string), &resSeguimiento); err == nil {
				seguimientos := make([]map[string]interface{}, 1)
				request.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)
				if fmt.Sprintf("%v", seguimientos) == "[]" {
					continue
				}

				var periodosSelecionados []map[string]interface{}
				for _, seguimiento := range seguimientos {
					for _, periodo := range periodos {
						if seguimiento["periodo_seguimiento_id"] == periodo["_id"] {
							for _, trimestre := range trimestres {
								var trimestreId float64
								if reflect.TypeOf(trimestre["Id"]).String() == "string" {
									trimestreId, _ = strconv.ParseFloat(trimestre["Id"].(string), 64)
								} else {
									trimestreId = trimestre["Id"].(float64)
								}
								var periodoId float64
								if reflect.TypeOf(periodo["periodo_id"]).String() == "string" {
									periodoId, _ = strconv.ParseFloat(periodo["periodo_id"].(string), 64)
								} else {
									periodoId = periodo["periodo_id"].(float64)
								}

								if trimestreId == periodoId {
									periodosSelecionados = append(periodosSelecionados, map[string]interface{}{"nombre": trimestre["ParametroId"].(map[string]interface{})["Nombre"].(string), "id": periodo["_id"]})
									break
								}
							}
							break
						}
					}
				}

				result = append(result, map[string]interface{}{"plan": plan["nombre"], "id": plan["_id"], "periodos": periodosSelecionados})
			} else {
				outputError = errors.New("error al decodificar el cuerpo de la solicitud")
			}
		}
	} else {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
	}
	return result, outputError
}
