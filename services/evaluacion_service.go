package services

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"

	"github.com/astaxie/beego"
	evaluacionhelper "github.com/udistrital/planeacion_evaluacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
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

func GetPlanesPeriodo(vigencia string, unidad string) (respuesta []map[string]interface{}, outputError error) {

	if len(vigencia) == 0 || len(unidad) == 0 {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
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
			outputError = errors.New("error al decodificar el cuerpo de la solicitud")
			return nil, outputError
		}
		periodos, err := evaluacionhelper.GetPeriodos(vigencia)
		if err != nil {
			return nil, errors.New("error GetPeriodos en la solicitud: " + err.Error())
		}
		trimestres, err := evaluacionhelper.GetTrimestres(vigencia)
		if err != nil {
			return nil, errors.New("error GetPeriodos en la solicitud: " + err.Error())
		}

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
				return nil, outputError
			}
		}
	} else {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud")
		return nil, outputError
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

	trimestres, err := evaluacionhelper.GetPeriodos(vigencia)
	if err != nil {
		return nil, errors.New("error GetPeriodos en la solicitud: " + err.Error())
	}

	if len(trimestres) == 0 {
		return nil, nil
	} else {
		i := 0
		for index, periodo := range trimestres {
			if periodo["_id"] == periodoId {
				i = index
				break
			}
		}
		evaluacion = evaluacionhelper.GetEvaluacion(plan, trimestres, i)
		return evaluacion, nil
	}
}

func Unidades(plan string, vigencia string) (interface{}, error) {
	if nombrePlan, err := url.QueryUnescape(plan); err == nil {
		if data, err := evaluacionhelper.GetUnidadesPorPlanYVigencia(nombrePlan, vigencia); err == nil {
			return data, nil
		} else {
			return nil, errors.New("Error obteniendo las unidades del plan y la vigencia dados ")
		}
	} else {
		return nil, errors.New("Error obteniendo las unidades del plan y la vigencia dados " + err.Error())

	}
}

func Avances(plan string, vigencia string, unidad string) (interface{}, error) {

	if nombrePlan, err1 := url.QueryUnescape(plan); err1 == nil {
		if data, err2 := evaluacionhelper.GetAvances(nombrePlan, vigencia, unidad); err2 == nil {
			return data, nil
		} else {
			return nil, errors.New("Error obteniendo los avances ")
		}
	} else {
		return nil, errors.New("Error obteniendo los avances " + err1.Error())
	}
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

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_AVALADO_PARA_SEGUIMIENTO, &respuestaEstado); err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
	}
	request.LimpiezaRespuestaRefactor(respuestaEstado, &estadoSeguimiento)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/tipo-seguimiento?query=activo:true,codigo_abreviacion:"+ABREVIACION_SEGUIMIENTO_PLAN_ACCION, &respuestaTipoSeguimiento); err != nil {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
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
				outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")
			}
		}
	} else {
		outputError = errors.New("error al decodificar el cuerpo de la solicitud: 404")

	}
	return planes, outputError
}
