package controllers

import (
	"net/http"
	"testing"
)

func TestGetPlanesPeriodo(t *testing.T) {
	if response, err := http.Get("http://localhost:8082/v1/evaluacion/planes_periodo/25/8"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetPlanesPeriodo Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetPlanesPeriodo Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetPlanesPeriodo: ", err.Error())
		t.Fail()
	}
}

func TestGetEvaluacion(t *testing.T) {
	if response, err := http.Get("http://localhost:8082/v1/evaluacion/25/63b5f7bb159830a9238fdbfd/635b1f995073f2675157dc7f"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetEvaluacion Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetEvaluacion Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetEvaluacion: ", err.Error())
		t.Fail()
	}
}
