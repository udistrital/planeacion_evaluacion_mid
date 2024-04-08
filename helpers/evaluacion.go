package helpers

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_mid/helpers"
	seguimientohelper "github.com/udistrital/planeacion_mid/helpers/seguimientoHelper"
	"github.com/udistrital/utils_oas/request"
)

var estadoHttp string = "500"

const (
	CodigoEstadoPlan        string = "A_SP" //6153355601c7a2365b2fb2a1
	CodigoEstadoSeguimiento string = "AV"   //622ba49216511e93a95c326d
	CodigoTipoSeguimiento   string = "S_SP" //61f236f525e40c582a0840d0
)

func GetEvaluacionTrimestre(planId string, periodoId string, actividadId string) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	actividades := make(map[string]interface{})
	detalle := make(map[string]interface{})

	idEstadoSeguimiento, err := getIdCodigoAbreviacion("estado-seguimiento", CodigoEstadoSeguimiento)
	if err != nil {
		panic(err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=estado_seguimiento_id:`+idEstadoSeguimiento+`,plan_id:`+planId+`,periodo_seguimiento_id:`+periodoId, &resSeguimiento); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &aux)
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
				helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
				detalle = seguimientohelper.ConvertirStringJson(detalle)
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

func GetEvaluacion(planId string, periodos []map[string]interface{}, trimestre int) []map[string]interface{} {
	var resSeguimiento map[string]interface{}
	var seguimiento map[string]interface{}
	var evaluacion []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	actividades := make(map[string]interface{})

	idEstadoSeguimiento, err := getIdCodigoAbreviacion("estado-seguimiento", CodigoEstadoSeguimiento)
	if err != nil {
		panic(err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=estado_seguimiento_id:`+idEstadoSeguimiento+`,plan_id:`+planId+`,periodo_seguimiento_id:`+periodos[trimestre]["_id"].(string), &resSeguimiento); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &aux)
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
					helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
					actividad = seguimientohelper.ConvertirStringJson(detalle)
				}
			} else {
				actividad = act.(map[string]interface{})
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
					if indexPeriodo == 0 {
						trimestreNom = "trimestre1"
					} else if indexPeriodo == 1 {
						trimestreNom = "trimestre2"
					} else if indexPeriodo == 2 {
						trimestreNom = "trimestre3"
					} else if indexPeriodo == 3 {
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

		helpers.SortSlice(&evaluacion, "numero")
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
				for _, trimestre := range []string{"trimestre1", "trimestre2", "trimestre3", "trimestre4"} {
					if val, ok := evaluacion[i][trimestre]; ok && fmt.Sprintf("%v", val) != "map[]" {
						meta := evaluacion[i][trimestre].(map[string]interface{})["meta"].(float64)
						if meta > 1 {
							switch trimestre {
							case "trimestre1":
								sum1 += 1.0
							case "trimestre2":
								sum2 += 1.0
							case "trimestre3":
								sum3 += 1.0
							case "trimestre4":
								sum4 += 1.0
							}
						} else {
							switch trimestre {
							case "trimestre1":
								sum1 += meta
							case "trimestre2":
								sum2 += meta
							case "trimestre3":
								sum3 += meta
							case "trimestre4":
								sum4 += meta
							}
						}
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
				updateActividad := func(trimestre string, cumplActividad interface{}) {
					key := fmt.Sprintf("%v", evaluacion[i][trimestre])
					if key != "map[]" {
						evaluacion[i][trimestre].(map[string]interface{})["actividad"] = cumplActividad
					}
				}

				updateActividad("trimestre1", cumplActividad1)
				updateActividad("trimestre2", cumplActividad2)
				updateActividad("trimestre3", cumplActividad3)
				updateActividad("trimestre4", cumplActividad4)
			}

		}

		return evaluacion
	}
	return nil
}

func GetPeriodos(vigencia string) []map[string]interface{} {
	var periodos []map[string]interface{}
	var resPeriodo map[string]interface{}
	var wg sync.WaitGroup
	trimestres := seguimientohelper.GetTrimestres(vigencia)
	periodosMutex := sync.Mutex{}

	for _, trimestre := range trimestres {
		wg.Add(1)
		if fmt.Sprintf("%v", trimestre) == "map[]" {
			wg.Done()
			continue
		}
		go func(trimestreId int, wg *sync.WaitGroup, periodos *[]map[string]interface{}) {
			periodosMutex.Lock()

			idTipoSeguimiento, err := getIdCodigoAbreviacion("tipo-seguimiento", CodigoTipoSeguimiento)
			if err != nil {
				panic(err.Error())
			}

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/periodo-seguimiento?fields=_id,periodo_id&query=tipo_seguimiento_id:`+idTipoSeguimiento+`,periodo_id:`+strconv.Itoa(trimestreId), &resPeriodo); err == nil {
				var periodo []map[string]interface{}
				helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
				(*periodos) = append((*periodos), periodo...)
			}
			periodosMutex.Unlock()
			wg.Done()
		}(int(trimestre["Id"].(float64)), &wg, &periodos)
	}

	wg.Wait()

	helpers.SortSlice(&periodos, "periodo_id")
	return periodos
}

func PlanDetalle(vigencia string, unidad string) (result []map[string]interface{}, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "PlanDetalle", "err": err, "status": estadoHttp}
			panic(outputError)
		}
	}()

	var resPlan map[string]interface{}
	var resSeguimiento map[string]interface{}

	idEstadoPlan, err := getIdCodigoAbreviacion("estado-plan", CodigoEstadoPlan)
	if err != nil {
		panic(err.Error())
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan?query=estado_plan_id:`+idEstadoPlan+`,dependencia_id:`+unidad+`,vigencia:`+vigencia, &resPlan); err == nil {
		planes := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(resPlan, &planes)
		if fmt.Sprintf("%v", planes) == "[]" {
			estadoHttp = "404"
			panic(err.Error())
		}

		periodos := GetPeriodos(vigencia)
		trimestres := seguimientohelper.GetTrimestres(vigencia)

		idEstadoSeguimiento, err1 := getIdCodigoAbreviacion("estado-seguimiento", CodigoEstadoSeguimiento)
		idTipoSeguimiento, err2 := getIdCodigoAbreviacion("tipo-seguimiento", CodigoTipoSeguimiento)
		if err1 != nil || err2 != nil {
			panic(err.Error())
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
				estadoHttp = "500"
				panic(err.Error())
			}
		}
	} else {
		estadoHttp = "500"
		panic(err.Error())
	}
	return result, outputError
}

func EvaluacionDetalle(vigencia string, plan string, periodoId string) (evaluacion []map[string]interface{}, outputError map[string]interface{}) {

	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{"function": "EvaluacionDetalle", "err": err, "status": estadoHttp}
			panic(outputError)
		}
	}()

	trimestres := GetPeriodos(vigencia)
	if len(trimestres) == 0 {
		estadoHttp = "401"
		panic("Error al obtener trimestres. " + "Error: " + estadoHttp)
	} else {
		i := 0
		for index, periodo := range trimestres {
			if periodo["_id"] == periodoId {
				i = index
				break
			}
		}

		evaluacion = GetEvaluacion(plan, trimestres, i)

	}
	return evaluacion, outputError
}

func getIdCodigoAbreviacion(ruta string, codigo string) (string, error) {
	var resEstado map[string]interface{}
	var estado []map[string]interface{}
	url := "http://" + beego.AppConfig.String("PlanesService") + "/" + ruta + "?query=activo:true,codigo_abreviacion:" + codigo
	err := request.GetJson(url, &resEstado)
	if err != nil {
		return "", err
	}
	request.LimpiezaRespuestaRefactor(resEstado, &estado)
	return estado[0]["_id"].(string), nil
}
