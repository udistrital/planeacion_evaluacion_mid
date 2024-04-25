package helpers

import (
	"encoding/json"

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
