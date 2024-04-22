package seguimientohelper

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"

	"log"
	"strings"

	"github.com/astaxie/beego"
	helpers "github.com/udistrital/planeacion_evaluacion_mid/helpers/utilidades"
	"github.com/udistrital/utils_oas/request"
)

func GetTrimestres(vigencia string) []map[string]interface{} {

	var res map[string]interface{}
	var trimestre []map[string]interface{}
	var trimestres []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T1", &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &trimestre)
		trimestres = append(trimestres, trimestre...)

		trimestre = nil
		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T2", &res); err == nil {
			helpers.LimpiezaRespuestaRefactor(res, &trimestre)
			trimestres = append(trimestres, trimestre...)

			trimestre = nil
			if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T3", &res); err == nil {
				helpers.LimpiezaRespuestaRefactor(res, &trimestre)
				trimestres = append(trimestres, trimestre...)

				trimestre = nil
				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=PeriodoId:"+vigencia+",ParametroId__CodigoAbreviacion:T4", &res); err == nil {
					helpers.LimpiezaRespuestaRefactor(res, &trimestre)
					trimestres = append(trimestres, trimestre...)
				} else {
					panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
				}
			} else {
				panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
			}
		} else {
			panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
		}
	} else {
		panic(map[string]interface{}{"funcion": "GetTrimestres", "err": "Error ", "status": "400", "log": err})
	}

	return trimestres
}

func GetActividades(subgrupo_id string) []map[string]interface{} {
	var res map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var datoPlan map[string]interface{}
	var actividades []map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle?query=subgrupo_id:"+subgrupo_id, &res); err == nil {
		aux := make([]map[string]interface{}, 1)
		helpers.LimpiezaRespuestaRefactor(res, &aux)
		subgrupoDetalle = aux[0]
		if subgrupoDetalle["dato_plan"] != nil {
			dato_plan_str := subgrupoDetalle["dato_plan"].(string)
			json.Unmarshal([]byte(dato_plan_str), &datoPlan)
			for indexActividad, element := range datoPlan {
				_ = indexActividad
				if err != nil {
					log.Panic(err)
				}
				if element.(map[string]interface{})["activo"] == true {
					actividades = append(actividades, element.(map[string]interface{}))
				}
			}
			sort.SliceStable(actividades, func(i, j int) bool {
				var a int
				var b int
				if reflect.TypeOf(actividades[j]["index"]).String() == "string" {
					b, _ = strconv.Atoi(actividades[j]["index"].(string))
				} else {
					b = int(actividades[j]["index"].(float64))
				}

				if reflect.TypeOf(actividades[i]["index"]).String() == "string" {
					a, _ = strconv.Atoi(actividades[i]["index"].(string))
				} else {
					a = int(actividades[i]["index"].(float64))
				}
				return a < b
			})
		}
	} else {
		panic(map[string]interface{}{"Code": "400", "Body": err, "Type": "error"})

	}
	return actividades
}

func GetActividad(seguimiento map[string]interface{}, index string, trimestre string) map[string]interface{} {
	var data map[string]interface{}
	var resEstado map[string]interface{}
	var resDetalle map[string]interface{}
	var informacion map[string]interface{}
	var cuantitativo map[string]interface{}
	cualitativo := map[string]interface{}{}
	evidencia := []interface{}{}
	evidenciaSeg := []map[string]interface{}{}
	estado := map[string]interface{}{}
	detalle := map[string]interface{}{}
	id := ""
	dato := make(map[string]interface{})
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[index] != nil {
		idS, segregado := dato[index].(map[string]interface{})["id"]

		if segregado && idS != "" {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[index].(map[string]interface{})["id"].(string), &resDetalle); err == nil {
				if resDetalle["Data"] != "null" {
					helpers.LimpiezaRespuestaRefactor(resDetalle, &detalle)
					detalle = ConvertirStringJson(detalle)

					id = detalle["_id"].(string)

					if len(detalle["informacion"].(map[string]interface{})) == 0 {
						informacion = GetInformacionPlan(seguimiento, index)
					} else {
						informacion = detalle["informacion"].(map[string]interface{})
					}

					if len(detalle["cuantitativo"].(map[string]interface{})) == 0 {
						cuantitativo = GetCuantitativoPlan(seguimiento, index, trimestre)
					} else {
						cuantitativo = detalle["cuantitativo"].(map[string]interface{})
					}

					if len(detalle["cualitativo"].(map[string]interface{})) == 0 {
						cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
					} else {
						cualitativo = detalle["cualitativo"].(map[string]interface{})
					}

					if len(detalle["evidencia"].([]map[string]interface{})) != 0 {
						evidenciaSeg = detalle["evidencia"].([]map[string]interface{})
					}

					if len(detalle["estado"].(map[string]interface{})) == 0 {
						if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &resEstado); err == nil {
							estado = map[string]interface{}{
								"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
								"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
							}
						}
					} else {
						estado = detalle["estado"].(map[string]interface{})
					}
				}
			}
		} else {
			if dato[index].(map[string]interface{})["informacion"] == nil {
				informacion = GetInformacionPlan(seguimiento, index)
			} else {
				informacion = dato[index].(map[string]interface{})["informacion"].(map[string]interface{})
			}

			if dato[index].(map[string]interface{})["cuantitativo"] == nil {
				cuantitativo = GetCuantitativoPlan(seguimiento, index, trimestre)
			} else {
				cuantitativo = dato[index].(map[string]interface{})["cuantitativo"].(map[string]interface{})
			}

			if dato[index].(map[string]interface{})["evidencia"] != nil {
				evidencia = dato[index].(map[string]interface{})["evidencia"].([]interface{})
			}
			if dato[index].(map[string]interface{})["estado"] == nil {
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &resEstado); err == nil {
					estado = map[string]interface{}{
						"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
						"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
					}
				}
			} else {
				estado = dato[index].(map[string]interface{})["estado"].(map[string]interface{})
			}

			if dato[index].(map[string]interface{})["cualitativo"] == nil {
				cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
			} else {
				cualitativo = dato[index].(map[string]interface{})["cualitativo"].(map[string]interface{})
			}
		}
	} else {
		if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:SRE", &resEstado); err == nil {
			estado = map[string]interface{}{
				"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
				"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
			}
		}
		informacion = GetInformacionPlan(seguimiento, index)
		cuantitativo = GetCuantitativoPlan(seguimiento, index, trimestre)
		cualitativo = map[string]interface{}{"reporte": "", "productos": "", "dificultades": ""}
	}

	data = map[string]interface{}{
		"id":           id,
		"informacion":  informacion,
		"cualitativo":  cualitativo,
		"cuantitativo": cuantitativo,
		"estado":       estado,
		"evidencia":    evidencia,
	}

	if id != "" {
		data["evidencia"] = evidenciaSeg
	}
	return data
}

func GetInformacionPlan(seguimiento map[string]interface{}, index string) map[string]interface{} {
	var resPlan map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var resInformacion map[string]interface{}
	var respuestaDependencia []map[string]interface{}
	var hijos []map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var periodo []map[string]interface{}

	informacion := map[string]interface{}{
		"ponderacion": "",
		"periodo":     "",
		"tarea":       "",
		"producto":    "",
		"nombre":      "",
		"descripcion": "",
		"index":       index,
		"unidad":      "",
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/plan/"+seguimiento["plan_id"].(string), &resPlan); err == nil {
		informacion["nombre"] = resPlan["Data"].(map[string]interface{})["nombre"]
		informacion["unidad"] = resPlan["Data"].(map[string]interface{})["dependencia_id"]
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

		if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
			helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
			informacion["trimestre"] = periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"]
		}
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+seguimiento["plan_id"].(string), &resInformacion); err == nil {
		helpers.LimpiezaRespuestaRefactor(resInformacion, &hijos)
		for _, hijo := range hijos {
			nombreHijo := strings.ToLower(hijo["nombre"].(string))
			if hijo["activo"] == true {
				var res map[string]interface{}

				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijo["_id"].(string), &res); err == nil {
					datoPlan := make(map[string]interface{})
					dato := make(map[string]interface{})

					nombreDetalle := strings.ToLower(res["Data"].([]interface{})[0].(map[string]interface{})["nombre"].(string))
					if strings.Contains(nombreDetalle, "indicadores") || strings.Contains(nombreDetalle, "indicador") {
						continue
					}

					json.Unmarshal([]byte(res["Data"].([]interface{})[0].(map[string]interface{})["dato"].(string)), &dato)
					if dato["required"] == false || dato["required"] == "false" {
						continue
					}

					json.Unmarshal([]byte(res["Data"].([]interface{})[0].(map[string]interface{})["dato_plan"].(string)), &datoPlan)

					if datoPlan[index] == nil {
						continue
					}

					switch {
					case strings.Contains(nombreHijo, "ponderación"):
						informacion["ponderacion"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "periodo") || strings.Contains(nombreHijo, "período"):
						informacion["periodo"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "tareas") || strings.Contains(nombreHijo, "actividades específicas"):
						informacion["tarea"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "producto"):
						informacion["producto"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					case strings.Contains(nombreHijo, "actividad general"):
						informacion["descripcion"] = datoPlan[index].(map[string]interface{})["dato"]
						continue
					}
				}
			}
		}
	}

	if err := request.GetJson("http://"+beego.AppConfig.String("OikosService")+"/dependencia_tipo_dependencia?query=DependenciaId:"+informacion["unidad"].(string), &respuestaDependencia); err == nil {
		informacion["unidad"] = respuestaDependencia[0]["DependenciaId"].(map[string]interface{})["Nombre"]
	} else {
		informacion["unidad"] = nil
	}

	return informacion
}

func GetCuantitativoPlan(seguimiento map[string]interface{}, index string, trimestre string) map[string]interface{} {
	var resInformacion map[string]interface{}
	var resDetalle map[string]interface{}
	var hijos []interface{}
	var subgrupos []map[string]interface{}
	var indicadores []map[string]interface{}
	var respuestas []map[string]interface{}
	response := map[string]interface{}{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/hijos/"+seguimiento["plan_id"].(string), &resInformacion); err == nil {
		helpers.LimpiezaRespuestaRefactor(resInformacion, &subgrupos)
		for _, subgrupo := range subgrupos {
			if strings.Contains(strings.ToLower(subgrupo["nombre"].(string)), "indicador") && subgrupo["activo"] == true {

				hijos = subgrupo["hijos"].([]interface{})
				hijos = append(hijos, subgrupo["_id"])

				for _, hijo := range hijos {
					var res map[string]interface{}
					if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo/"+hijo.(string), &res); err == nil {
						hijosIndicadores := res["Data"].(map[string]interface{})["hijos"].([]interface{})

						var dato_plan map[string]interface{}

						informacion := map[string]interface{}{
							"detalleReporte": "",
						}

						respuesta := map[string]interface{}{
							"indicador":            0,
							"indicadorAcumulado":   0,
							"avanceAcumulado":      0,
							"brechaExistente":      0,
							"acumuladoNumerador":   0,
							"acumuladoDenominador": 0,
							"meta":                 0,
						}

						for _, hijoI := range hijosIndicadores {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+hijoI.(string), &resDetalle); err == nil {
								var subgrupo_detalle []map[string]interface{}
								helpers.LimpiezaRespuestaRefactor(resDetalle, &subgrupo_detalle)

								if len(subgrupo_detalle) > 0 {
									if subgrupo_detalle[0]["dato_plan"] != nil {
										dato_plan_str := subgrupo_detalle[0]["dato_plan"].(string)
										json.Unmarshal([]byte(dato_plan_str), &dato_plan)
										nombreDetalle := strings.ToLower(subgrupo_detalle[0]["nombre"].(string))
										if dato_plan[index] == nil || dato_plan[index].(map[string]interface{})["dato"] == "" {
											break
										}

										switch {
										case strings.Contains(nombreDetalle, "nombre"):
											informacion["nombre"] = dato_plan[index].(map[string]interface{})["dato"]
											respuesta["nombre"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "meta"):
											informacion["meta"] = dato_plan[index].(map[string]interface{})["dato"]
											if reflect.TypeOf(dato_plan[index].(map[string]interface{})["dato"]).String() == "string" {
												respuesta["meta"], _ = strconv.ParseFloat(dato_plan[index].(map[string]interface{})["dato"].(string), 64)
											} else {
												respuesta["meta"] = dato_plan[index].(map[string]interface{})["dato"].(float64)
											}
											continue
										case strings.Contains(nombreDetalle, "fórmula"):
											informacion["formula"] = dato_plan[index].(map[string]interface{})["dato"]
											continue
										case strings.Contains(nombreDetalle, "criterio"):
											informacion["denominador"] = dato_plan[index].(map[string]interface{})["dato"]
											if informacion["denominador"] == "Denominador fijo" {
												// informacion["reporteDenominador"] = GetDenominadorFijo(seguimiento, len(indicadores), index)
											}
											continue
										case strings.Contains(nombreDetalle, "tendencia"):
											informacion["tendencia"] = strings.Trim(dato_plan[index].(map[string]interface{})["dato"].(string), " ")
											continue
										case strings.Contains(nombreDetalle, "unidad de medida"):
											informacion["unidad"] = strings.Trim(dato_plan[index].(map[string]interface{})["dato"].(string), " ")
											respuesta["unidad"] = strings.Trim(dato_plan[index].(map[string]interface{})["dato"].(string), " ")
											continue
										}
									}
								}
							}
						}

						if informacion["reporteDenominador"] == 1.0 {
							informacion["reporteDenominador"] = nil
						}

						if informacion["nombre"] != nil && informacion["nombre"] != "" {
							indicadores = append(indicadores, informacion)
							respuestas = append(respuestas, respuesta)
						}

						if informacion["denominador"] == nil {
							informacion["denominador"] = ""
						}

						respuestas = GetRespuestaAnterior(seguimiento, len(indicadores)-1, respuestas, index, trimestre)
					}
				}
				break
			}
		}
	}

	response["indicadores"] = indicadores
	response["resultados"] = respuestas
	return response
}

func GetDenominadorFijo(dataSeg map[string]interface{}, index int, indexActividad string) float64 {
	plan_id := dataSeg["plan_id"].(string)
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var seguimientos []map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var periodo []map[string]interface{}

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_id, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)

		for _, seguimiento := range seguimientos {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
					helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
					trimestre := periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"]
					if trimestre == "T1" {
						if seguimiento["dato"] != "{}" {
							dato := make(map[string]interface{})
							datoStr := seguimiento["dato"].(string)
							json.Unmarshal([]byte(datoStr), &dato)
							if dato[indexActividad] == nil {
								break
							}

							seguimientoActividad := dato[indexActividad].(map[string]interface{})
							if seguimientoActividad["cuantitativo"] == nil {
								break
							}

							if fmt.Sprint(reflect.TypeOf(seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"])) == "int" || fmt.Sprint(reflect.TypeOf(seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"])) == "float64" {
								return seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"].(float64)
							} else {
								aux2, err := strconv.ParseFloat(seguimientoActividad["cuantitativo"].(map[string]interface{})["indicadores"].([]interface{})[index].(map[string]interface{})["reporteDenominador"].(string), 64)
								if err == nil {
									return aux2
								}
							}
						}
						break
					}
				}
			}
		}
	}

	return 1
}

func GetRespuestaAnterior(dataSeg map[string]interface{}, index int, respuestas []map[string]interface{}, indexActividad string, trimestre string) []map[string]interface{} {
	plan_id := dataSeg["plan_id"].(string)
	var resSeguimiento map[string]interface{}
	var resPeriodoSeguimiento map[string]interface{}
	var resPeriodo map[string]interface{}
	var periodoSeguimiento map[string]interface{}
	var seguimientos []map[string]interface{}
	var periodo []map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento?query=activo:true,plan_id:"+plan_id, &resSeguimiento); err == nil {
		helpers.LimpiezaRespuestaRefactor(resSeguimiento, &seguimientos)

		acumuladoNumerador := 0.0
		acumuladoDenominador := 0.0
		indicadorAcumulado := 0.0
		avanceAcumulado := 0.0
		brechaExistente := 0.0
		divisionCero := false
		for _, seguimiento := range seguimientos {

			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/periodo-seguimiento/"+seguimiento["periodo_seguimiento_id"].(string), &resPeriodoSeguimiento); err == nil {
				helpers.LimpiezaRespuestaRefactor(resPeriodoSeguimiento, &periodoSeguimiento)

				if err := request.GetJson("http://"+beego.AppConfig.String("ParametrosService")+"/parametro_periodo?query=Id:"+periodoSeguimiento["periodo_id"].(string), &resPeriodo); err == nil {
					helpers.LimpiezaRespuestaRefactor(resPeriodo, &periodo)
					tri, _ := strconv.Atoi(string(trimestre[1]))
					segTrimestre, _ := strconv.Atoi(string(periodo[0]["ParametroId"].(map[string]interface{})["CodigoAbreviacion"].(string)[1]))
					if (tri - 1) == segTrimestre {
						if seguimiento["dato"] != "{}" {
							dato := make(map[string]interface{})
							datoStr := seguimiento["dato"].(string)
							json.Unmarshal([]byte(datoStr), &dato)
							if dato[indexActividad] == nil {
								respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
								respuestas[index]["avanceAcumulado"] = avanceAcumulado
								respuestas[index]["brechaExistente"] = brechaExistente
								respuestas[index]["divisionCero"] = divisionCero
								continue
							}

							id, segregado := dato[indexActividad].(map[string]interface{})["id"]
							if segregado && id != "" {
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
									helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
									detalle = ConvertirStringJson(detalle)

									if fmt.Sprintf("%v", detalle["cuantitativo"]) == "map[]" {
										respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
										respuestas[index]["avanceAcumulado"] = avanceAcumulado
										respuestas[index]["brechaExistente"] = brechaExistente
										respuestas[index]["divisionCero"] = divisionCero
										continue
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"] != nil {
										indicadorAcumulado += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"] != nil {
										avanceAcumulado += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"] != nil {
										brechaExistente += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"] != nil {
										divisionCero = detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"].(bool)
									} else {
										divisionCero = false
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"] != nil {
										acumuladoDenominador += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"].(float64)
									}

									if detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"] != nil {
										acumuladoNumerador += detalle["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"].(float64)
									}
								}
							} else {
								seguimientoActividad := dato[indexActividad].(map[string]interface{})
								if seguimientoActividad["cuantitativo"] == nil {
									respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
									respuestas[index]["avanceAcumulado"] = avanceAcumulado
									respuestas[index]["brechaExistente"] = brechaExistente
									respuestas[index]["divisionCero"] = divisionCero
									continue
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"] != nil {
									indicadorAcumulado += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["indicadorAcumulado"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"] != nil {
									avanceAcumulado += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["avanceAcumulado"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"] != nil {
									brechaExistente += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["brechaExistente"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"] != nil {
									divisionCero = seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["divisionCero"].(bool)
								} else {
									divisionCero = false
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"] != nil {
									acumuladoDenominador += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoDenominador"].(float64)
								}

								if seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"] != nil {
									acumuladoNumerador += seguimientoActividad["cuantitativo"].(map[string]interface{})["resultados"].([]interface{})[index].(map[string]interface{})["acumuladoNumerador"].(float64)
								}
							}
						}

						respuestas[index]["indicadorAcumulado"] = indicadorAcumulado
						respuestas[index]["avanceAcumulado"] = avanceAcumulado
						respuestas[index]["brechaExistente"] = brechaExistente
						respuestas[index]["acumuladoNumerador"] = acumuladoNumerador
						respuestas[index]["acumuladoDenominador"] = acumuladoDenominador
						respuestas[index]["divisionCero"] = divisionCero
						break
					}
				}
			}

		}
	}

	return respuestas
}

func GetSubgrupoDetalle(subgrupo_id string, index string) map[string]interface{} {
	var respuesta map[string]interface{}
	var subgrupoDetalle map[string]interface{}
	var datoPlan map[string]interface{}
	var data map[string]interface{}
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo-detalle/detalle/"+subgrupo_id, &respuesta); err != nil {
		panic(map[string]interface{}{"funcion": "GuardarPlan", "err": "Error get subgrupo-detalle \"key\"", "status": "400", "log": err})
	}
	aux := make([]map[string]interface{}, 1)
	helpers.LimpiezaRespuestaRefactor(respuesta, &aux)
	subgrupoDetalle = aux[0]
	if subgrupoDetalle["dato_plan"] != nil {
		dato_plan_str := subgrupoDetalle["dato_plan"].(string)
		json.Unmarshal([]byte(dato_plan_str), &datoPlan)

		data = datoPlan[index].(map[string]interface{})
	}
	return data
}

func GetEstadoSeguimiento(seguimiento map[string]interface{}) string {
	var resEstado map[string]interface{}
	enReporte := true
	estado := map[string]interface{}{}
	dato := make(map[string]interface{})

	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento/"+seguimiento["estado_seguimiento_id"].(string), &resEstado); err == nil {
		estado = map[string]interface{}{
			"nombre": resEstado["Data"].(map[string]interface{})["nombre"],
			"id":     resEstado["Data"].(map[string]interface{})["_id"],
		}

		for _, actividad := range dato {
			_, datosUnidos := actividad.(map[string]interface{})["estado"]
			if datosUnidos {
				if actividad.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad en reporte" {
					enReporte = false
				}
			} else {
				var resSeguimientoDetalle map[string]interface{}
				var seguimientoDetalle map[string]interface{}
				if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+actividad.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
					helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &seguimientoDetalle)
					dato := make(map[string]interface{})
					json.Unmarshal([]byte(seguimientoDetalle["estado"].(string)), &dato)
					if dato["nombre"] != "Actividad en reporte" {
						enReporte = false
					}
				}
			}
		}

		if enReporte {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/estado-seguimiento?query=codigo_abreviacion:ER", &resEstado); err == nil {
				estado = map[string]interface{}{
					"nombre": resEstado["Data"].([]interface{})[0].(map[string]interface{})["nombre"],
					"id":     resEstado["Data"].([]interface{})[0].(map[string]interface{})["_id"],
				}
			}
		}
	}

	return estado["id"].(string)
}

func ActividadReportable(seguimiento map[string]interface{}, indexActividad string) (bool, map[string]interface{}) {
	dato := make(map[string]interface{})
	estado := map[string]interface{}{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := map[string]interface{}{}
	var cuantitativo interface{}
	var cualitativo interface{}
	datoStr := seguimiento["dato"].(string)
	json.Unmarshal([]byte(datoStr), &dato)

	if dato[indexActividad] == nil {
		return false, map[string]interface{}{"error": 1, "motivo": "Actividad sin seguimiento"}
	} else {
		_, datosUnidos := dato[indexActividad].(map[string]interface{})["estado"]
		if datosUnidos {
			estado = dato[indexActividad].(map[string]interface{})["estado"].(map[string]interface{})
			cuantitativo = dato[indexActividad].(map[string]interface{})["cuantitativo"]
			cualitativo = dato[indexActividad].(map[string]interface{})["cualitativo"]
		} else {
			if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+dato[indexActividad].(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
				helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
				detalle = ConvertirStringJson(detalle)
				estado = detalle["estado"].(map[string]interface{})
				cualitativo = detalle["cualitativo"]
				cuantitativo = detalle["cuantitativo"]
			}

			if fmt.Sprintf("%v", cuantitativo) == "map[]" {
				cuantitativo = nil
			}

			if fmt.Sprintf("%v", cualitativo) == "map[]" {
				cualitativo = nil
			}
		}

		if estado["nombre"] != "Actividad en reporte" {
			return false, map[string]interface{}{"error": 2, "motivo": "El estado de la actividad no es el adecuado"}
		}

		if cuantitativo == nil {
			return false, map[string]interface{}{"error": 3, "motivo": "Componenten cuantitativo sin guardar"}
		}

		if cualitativo == nil {
			return false, map[string]interface{}{"error": 4, "motivo": "Componenten cualitativo sin guardar"}
		} else {
			cualitativo := cualitativo.(map[string]interface{})
			if cualitativo["dificultades"] == "" || cualitativo["productos"] == "" || cualitativo["reporte"] == "" {
				return false, map[string]interface{}{"error": 5, "motivo": "Campos vacios en el componenten cualitativo"}
			}
		}

	}

	return true, nil
}

func SeguimientoReportable(seguimiento map[string]interface{}) (bool, map[string]interface{}) {
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {

				actividades := GetActividades(subgrupos[i]["_id"].(string))
				if seguimiento["dato"] == "{}" {
					for _, actividad := range actividades {
						dato[actividad["index"].(string)] = actividad["dato"]
					}
					return false, map[string]interface{}{"error": 1, "motivo": "No hay actividades resportadas", "actividades": dato}
				} else {
					dato_plan_str := seguimiento["dato"].(string)
					json.Unmarshal([]byte(dato_plan_str), &datoPlan)

					for indexActividad, element := range datoPlan {
						id, segregado := element.(map[string]interface{})["id"]
						if segregado && id != "" {
							if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+id.(string), &resSeguimientoDetalle); err == nil {
								helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
								detalle = ConvertirStringJson(detalle)
							}

							for _, actividad := range actividades {
								if reflect.TypeOf(actividad["index"]).String() == "string" {
									if indexActividad == actividad["index"] {
										actividad["estado"] = detalle["estado"]
									}
								} else {
									if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
										actividad["estado"] = detalle["estado"]
									}
								}
							}
						} else {
							for _, actividad := range actividades {
								if reflect.TypeOf(actividad["index"]).String() == "string" {
									if indexActividad == actividad["index"] {
										actividad["estado"] = element.(map[string]interface{})["estado"]
									}
								} else {
									if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
										actividad["estado"] = element.(map[string]interface{})["estado"]
									}
								}
							}
						}
					}

					for _, actividad := range actividades {
						if actividad["estado"] == nil {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								dato[actividad["index"].(string)] = actividad["dato"]
							} else {
								dato[strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)] = actividad["dato"]
							}
						} else if actividad["estado"].(map[string]interface{})["nombre"] != "Actividad reportada" && actividad["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" {
							if reflect.TypeOf(actividad["index"]).String() == "string" {
								dato[actividad["index"].(string)] = actividad["dato"]
							} else {
								dato[strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64)] = actividad["dato"]
							}
						}
					}

					if fmt.Sprintf("%v", dato) != "map[]" {
						return false, map[string]interface{}{"error": 2, "motivo": "Hay actividades sin resportar", "actividades": dato}
					} else {
						return true, nil
					}
				}
			}
		}
	}
	return true, nil
}

func ActividadConObservaciones(seguimiento map[string]interface{}) bool {
	var cuantitativo map[string]interface{}
	var cualitativo map[string]interface{}

	if seguimiento["cuantitativo"] != nil {
		cuantitativo = seguimiento["cuantitativo"].(map[string]interface{})
		for _, indicador := range cuantitativo["indicadores"].([]interface{}) {
			if indicador.(map[string]interface{})["observaciones"] != "" && indicador.(map[string]interface{})["observaciones"] != "Sin observación" && indicador.(map[string]interface{})["observaciones"] != nil {
				return true
			}
		}
	}

	if seguimiento["cualitativo"] != nil {
		cualitativo = seguimiento["cualitativo"].(map[string]interface{})
		if cualitativo["observaciones"] != "" && cualitativo["observaciones"] != "Sin observación" && cualitativo["observaciones"] != nil {
			return true
		}
	}

	if seguimiento["evidencia"] != nil {
		for _, evidencia := range seguimiento["evidencia"].([]map[string]interface{}) {
			if evidencia["Observacion"] != "" && evidencia["Observacion"] != "Sin observación" && evidencia["Observacion"] != nil {
				return true
			}
		}
	}

	return false
}

func SeguimientoAvalable(seguimiento map[string]interface{}) (bool, bool, map[string]interface{}) {
	var res map[string]interface{}
	var subgrupos []map[string]interface{}
	var datoPlan map[string]interface{}
	var resSeguimientoDetalle map[string]interface{}
	detalle := make(map[string]interface{})
	dato := make(map[string]interface{})
	observaciones := false
	avaladas := false
	estado := map[string]interface{}{}

	planId := seguimiento["plan_id"].(string)
	if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/subgrupo?query=padre:"+planId, &res); err == nil {
		helpers.LimpiezaRespuestaRefactor(res, &subgrupos)

		for i := 0; i < len(subgrupos); i++ {
			if strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "actividad") && strings.Contains(strings.ToLower(subgrupos[i]["nombre"].(string)), "general") {

				actividades := GetActividades(subgrupos[i]["_id"].(string))

				dato_plan_str := seguimiento["dato"].(string)
				json.Unmarshal([]byte(dato_plan_str), &datoPlan)

				for indexActividad, element := range datoPlan {
					id, segregado := element.(map[string]interface{})["id"]

					for _, actividad := range actividades {
						if reflect.TypeOf(actividad["index"]).String() == "string" {
							if indexActividad == actividad["index"] {
								if segregado && id != "" {
									if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
										helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
										detalle = ConvertirStringJson(detalle)
										estado = detalle["estado"].(map[string]interface{})
										if estado["nombre"] != "Actividad avalada" && estado["nombre"] != "Con observaciones" {
											dato[indexActividad] = actividad["dato"]
										}
									}
								} else {
									if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
										dato[indexActividad] = actividad["dato"]
									}
								}
							}
						} else if indexActividad == strconv.FormatFloat(actividad["index"].(float64), 'g', 5, 64) {
							if segregado && id != "" {
								if err := request.GetJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+element.(map[string]interface{})["id"].(string), &resSeguimientoDetalle); err == nil {
									helpers.LimpiezaRespuestaRefactor(resSeguimientoDetalle, &detalle)
									detalle = ConvertirStringJson(detalle)
									estado = detalle["estado"].(map[string]interface{})
									if estado["nombre"] != "Actividad avalada" && estado["nombre"] != "Con observaciones" {
										dato[indexActividad] = actividad["dato"]
									}
								}
							} else {
								if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Actividad avalada" && element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] != "Con observaciones" {
									dato[indexActividad] = actividad["dato"]
								}
							}
						}
					}

					if segregado && id != "" {
						if estado["nombre"] == "Con observaciones" {
							observaciones = true
						}

						if estado["nombre"] == "Actividad avalada" {
							avaladas = true
						}
					} else {
						if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Con observaciones" {
							observaciones = true
						}

						if element.(map[string]interface{})["estado"].(map[string]interface{})["nombre"] == "Actividad avalada" {
							avaladas = true
						}
					}
				}
			}
		}
	}

	if fmt.Sprintf("%v", dato) != "map[]" {
		return false, false, map[string]interface{}{"error": 1, "motivo": "Hay actividades sin revisar", "actividades": dato}
	}

	return avaladas, observaciones, nil
}

func ConvertirJsonString(diccionario map[string]interface{}) map[string]interface{} {
	dicStrings := map[string]interface{}{}

	for clave, valor := range diccionario {
		jsonBytes, _ := json.Marshal(valor)
		if string(jsonBytes) == "null" {
			dicStrings[clave] = "{}"
		} else if clave == "_id" || clave == "fecha_creacion" {
			dicStrings[clave] = valor
		} else {
			dicStrings[clave] = string(jsonBytes)
		}
	}

	return dicStrings
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

func GuardarDetalleSegimiento(detalle map[string]interface{}, actualizar bool) string {
	var res map[string]interface{}
	var id string

	detalle = ConvertirJsonString(detalle)

	if _, existe := detalle["informacion"]; !existe {
		detalle["informacion"] = "{}"
	}
	if _, existe := detalle["cualitativo"]; !existe {
		detalle["cualitativo"] = "{}"
	}
	if _, existe := detalle["cuantitativo"]; !existe {
		detalle["cuantitativo"] = "{}"
	}
	if _, existe := detalle["evidencia"]; !existe {
		detalle["evidencia"] = "[]"
	}

	if actualizar {
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle/"+detalle["_id"].(string), "PUT", &res, detalle); err == nil {
			aux := make(map[string]interface{})
			helpers.LimpiezaRespuestaRefactor(res, &aux)
			id = aux["_id"].(string)
		}
	} else {
		if err := helpers.SendJson("http://"+beego.AppConfig.String("PlanesService")+"/seguimiento-detalle", "POST", &res, detalle); err == nil {
			aux := make(map[string]interface{})
			helpers.LimpiezaRespuestaRefactor(res, &aux)
			id = aux["_id"].(string)
		}
	}

	return id
}
