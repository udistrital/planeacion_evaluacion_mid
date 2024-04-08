package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_evaluacion_mid/helpers"
	"github.com/udistrital/utils_oas/request"
)

// EvaluacionController operations for Evaluacion
type EvaluacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *EvaluacionController) URLMapping() {
	c.Mapping("GetEvaluacion", c.GetEvaluacion)
	c.Mapping("GetPlanesPeriodo", c.GetPlanesPeriodo)
}

// PlanesPeriodo ...
// @Title PlanesPeriodo
// @Description get Planes y vigencias para la unidad y vigencia dado
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	unidad 		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /planes_periodo/:vigencia/:unidad [get]
func (c *EvaluacionController) GetPlanesPeriodo() {

	defer request.ErrorController(c.Controller, "EvaluacionController")

	vigencia := c.Ctx.Input.Param(":vigencia")
	unidad := c.Ctx.Input.Param(":unidad")

	if len(vigencia) == 0 || len(unidad) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request contains incorrect params", "Data": nil}
		c.ServeJSON()
		return
	}

	if data, err := helpers.PlanDetalle(vigencia, unidad); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
	} else {
		panic(map[string]interface{}{"Success": false, "Status": "400", "Message": "Error in GetPlanesPeriodo", "Data": nil, "Error": err})
	}

	c.ServeJSON()
}

// Evaluacion ...
// @Title Evaluacion
// @Description get Evaluacion
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	plan 		path 	string	true		"The key for staticblock"
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:vigencia/:plan:/:periodo [get]
func (c *EvaluacionController) GetEvaluacion() {

	defer request.ErrorController(c.Controller, "EvaluacionController")

	vigencia := c.Ctx.Input.Param(":vigencia")
	plan := c.Ctx.Input.Param(":plan")
	periodoId := c.Ctx.Input.Param(":periodo")

	if len(vigencia) == 0 || len(plan) == 0 || len(periodoId) == 0 {
		c.Data["json"] = map[string]interface{}{"Success": false, "Status": "404", "Message": "Request containt incorrect params", "Data": nil}
		c.ServeJSON()
		return
	}

	if data, err := helpers.EvaluacionDetalle(vigencia, plan, periodoId); err == nil {
		c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Successful", "Data": data}
	} else {
		panic(map[string]interface{}{"Success": false, "Status": "400", "Message": "Error in GetEvaluacion", "Data": nil, "Error": err})
	}

	c.ServeJSON()
}
