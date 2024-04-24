package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/astaxie/beego"
	"github.com/udistrital/utils_oas/request"
)

var estadoHttp string = "500"

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

func ConvertirStringJson(diccionario map[string]interface{}) map[string]interface{} {
	dicStrings := map[string]interface{}{}
	for clave, valor := range diccionario {
		if clave == "informacion" || clave == "cualitativo" || clave == "cuantitativo" || clave == "estado" {
			datoJson := make(map[string]interface{})
			json.Unmarshal([]byte(valor.(string)), &datoJson)
			dicStrings[clave] = datoJson
		} else if clave == "evidencia" {
			var datoJson []map[string]interface{}
			json.Unmarshal([]byte(valor.(string)), &datoJson)
			dicStrings[clave] = datoJson
		} else {
			dicStrings[clave] = valor
		}
	}
	return dicStrings
}

func GetTrimestres(vigencia string) ([]map[string]interface{}, error) {

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
					return nil, errors.New("error GetTrimestres en la solicitud: " + err.Error())
				}
			} else {
				return nil, errors.New("error GetTrimestres en la solicitud: " + err.Error())
			}
		} else {
			return nil, errors.New("error GetTrimestres en la solicitud: " + err.Error())
		}
	} else {
		return nil, errors.New("error GetTrimestres en la solicitud: " + err.Error())
	}
	return trimestres, nil
}

func GetPeriodos(vigencia string) ([]map[string]interface{}, error) {
	var periodos []map[string]interface{}
	var resPeriodo map[string]interface{}
	var wg sync.WaitGroup
	trimestres, err := GetTrimestres(vigencia)
	if err != nil {
		return nil, errors.New("error GetTrimestres en la solicitud: " + err.Error())
	}

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
				(*periodos) = append((*periodos), periodo...)
			} else {
				err = errors.New("error al decodificar el cuerpo de la solicitud")
			}
			periodosMutex.Unlock()
			wg.Done()
		}(int(trimestre["Id"].(float64)), &wg, &periodos)
	}

	wg.Wait()

	request.SortSlice(&periodos, "periodo_id")
	return periodos, nil
}

func GetIdCodigoAbreviacion(ruta string, codigo string) (string, error) {
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

func FiltrarArreglo(data []map[string]interface{}, condicion func(map[string]interface{}) bool) []map[string]interface{} {
	fltd := make([]map[string]interface{}, 0)
	for _, v := range data {
		if condicion(v) {
			fltd = append(fltd, v)
		}
	}
	return fltd
}

func PlanDetalle(vigencia string, unidad string) (result []map[string]interface{}, outputError error) {
	defer func() {
		if err := recover(); err != nil {
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
			estadoHttp = "404"
		}
	}()

	var resPlan map[string]interface{}
	var resSeguimiento map[string]interface{}

	idEstadoPlan, err := GetIdCodigoAbreviacion("estado-plan", CodigoEstadoPlan)
	if err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/plan?query=estado_plan_id:`+idEstadoPlan+`,dependencia_id:`+unidad+`,vigencia:`+vigencia, &resPlan); err == nil {
		planes := make([]map[string]interface{}, 1)
		request.LimpiezaRespuestaRefactor(resPlan, &planes)
		if fmt.Sprintf("%v", planes) == "[]" {
			estadoHttp = "404"
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		}

		periodos, err := GetPeriodos(vigencia)
		if err != nil {
			return nil, errors.New("error GetPeriodos en la solicitud: " + err.Error())
		}
		trimestres, err := GetTrimestres(vigencia)
		if err != nil {
			return nil, errors.New("error GetPeriodos en la solicitud: " + err.Error())
		}

		idEstadoSeguimiento, err1 := GetIdCodigoAbreviacion("estado-seguimiento", CodigoEstadoSeguimiento)
		idTipoSeguimiento, err2 := GetIdCodigoAbreviacion("tipo-seguimiento", CodigoTipoSeguimiento)
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
				estadoHttp = "500"
				outputError = errors.New("error al decodificar el cuerpo de la solicitud")
			}
		}
	} else {
		estadoHttp = "500"
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
	}
	return result, outputError
}

func GetPlanesParaEvaluar() (planes []string, outputError map[string]interface{}) {
	defer func() {
		if err := recover(); err != nil {
			outputError = map[string]interface{}{
				"funcion": "GetPlanesParaEvaluar",
				"err":     err,
				"status":  "400",
			}
			estadoHttp = "404"
			outputError = map[string]interface{}{
				"err":    outputError,
				"status": estadoHttp,
			}
		}
	}()

	var respuestaEstado map[string]interface{}
	var respuestaTipoSeguimiento map[string]interface{}
	var respuestaSeguimiento map[string]interface{}

	var estadoSeguimiento []map[string]interface{}
	var tipoSeguimiento []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_AVALADO_PARA_SEGUIMIENTO, &respuestaEstado); err != nil {
		estadoHttp = "404"
		outputError = map[string]interface{}{
			"err":    outputError,
			"status": estadoHttp,
		}
	}
	request.LimpiezaRespuestaRefactor(respuestaEstado, &estadoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_SEGUIMIENTO_PLAN_ACCION, &respuestaTipoSeguimiento); err != nil {
		estadoHttp = "404"
		outputError = map[string]interface{}{
			"err":    outputError,
			"status": estadoHttp,
		}
	}
	request.LimpiezaRespuestaRefactor(respuestaTipoSeguimiento, &tipoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+`/seguimiento?query=tipo_seguimiento_id:`+tipoSeguimiento[0]["_id"].(string)+`,estado_seguimiento_id:`+estadoSeguimiento[0]["_id"].(string), &respuestaSeguimiento); err == nil {
		var seguimientos []map[string]interface{}
		request.LimpiezaRespuestaRefactor(respuestaSeguimiento, &seguimientos)
		for _, seguimiento := range seguimientos {
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
					planes = append(planes, plan["nombre"].(string))
				}
			} else {
				estadoHttp = "404"
				outputError = map[string]interface{}{
					"err":    outputError,
					"status": estadoHttp,
				}
			}
		}
	} else {
		estadoHttp = "404"
		outputError = map[string]interface{}{
			"err":    outputError,
			"status": estadoHttp,
		}
	}
	return planes, outputError
}
