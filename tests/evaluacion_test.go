package controllers

import (
	"net/http"
	"testing"
)

// SE NECESITAN DATOS PARA PODER VALIDAR EL CASO
func TestGetPlanesPeriodo(t *testing.T) {
	if response, err := http.Get("http://localhost:8082/v1/evaluacion/planes-periodo/25/8"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestGetPlanesPeriodo Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestGetPlanesPeriodo Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestGetPlanesPeriodo:", err.Error())
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
		t.Error("Error TestGetEvaluacion:", err.Error())
		t.Fail()
	}
}

func TestPlanesAEvaluar(t *testing.T) {

	if response, err := http.Get("http://localhost:8082/v1/evaluacion/planes"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestPlanesAEvaluar Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPlanesAEvaluar Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPlanesAEvaluar:", err.Error())
		t.Fail()
	}
}
func TestUnidades(t *testing.T) {

	if response, err := http.Get("http://localhost:8082/v1/evaluacion/unidades/Plan%20de%20acci%C3%B3n%202023%20Prod%20Seguimiento/25"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestUnidades Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestUnidades Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestUnidades:", err.Error())
		t.Fail()
	}
}

func TestAvance(t *testing.T) {

	if response, err := http.Get("http://localhost:8082/v1/evaluacion/avance/Plan%20de%20acci%C3%B3n%202023%20Prod%20Seguimiento/25/14"); err == nil {
		if response.StatusCode != 200 {
			t.Error("Error TestAvance Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestAvance Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestAvance:", err.Error())
		t.Fail()
	}
}
