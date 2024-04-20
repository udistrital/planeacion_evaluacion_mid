package services

import (
	"errors"
	"fmt"
	"net/url"

	evaluacionhelper "github.com/udistrital/planeacion_evaluacion_mid/helpers"
)

func GetPlanesPeriodo(vigencia string, unidad string) (interface{}, error) {

	if len(vigencia) == 0 || len(unidad) == 0 {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: ")
	}
	if respuesta, err := evaluacionhelper.GetPlanesPeriodo(unidad, vigencia); err == nil {
		return respuesta, nil
	} else {
		return nil, errors.New("error al decodificar el cuerpo de la solicitud: ")
	}
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

	trimestres := evaluacionhelper.GetPeriodos(vigencia)
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
	// Imprimir los parámetros recibidos
	fmt.Println("Plan recibido en service:", plan)
	fmt.Println("Vigencia recibida en service:", vigencia)
	if nombrePlan, err := url.QueryUnescape(plan); err == nil {
		if data, err := evaluacionhelper.GetUnidadesPorPlanYVigencia(nombrePlan, vigencia); err == nil {
			fmt.Println("Plan recibido en service dentro de la funcion GetUnidadesPorPlanYVigencia:", plan)
			fmt.Println("Plan recibido en service dentro de la funcion GetUnidadesPorPlanYVigencia:", nombrePlan)
			fmt.Println("Vigencia recibida en service dentro de la funcion GetUnidadesPorPlanYVigencia:", vigencia)
			return data, nil
		} else {
			fmt.Println("Plan recibido dentro del else:", plan)
			fmt.Println("Plan recibido dentro del else:", nombrePlan)
			fmt.Println("Vigencia recibida dentro del else:", vigencia)
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
