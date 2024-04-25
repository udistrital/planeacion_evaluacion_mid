package controllers

import (
	"github.com/astaxie/beego"
	"github.com/udistrital/planeacion_evaluacion_mid/services"
	"github.com/udistrital/utils_oas/errorhandler"
	"github.com/udistrital/utils_oas/requestresponse"
)

// EvaluacionController operations for Evaluacion
type EvaluacionController struct {
	beego.Controller
}

// URLMapping ...
func (c *EvaluacionController) URLMapping() {
	c.Mapping("GetPlanesPeriodo", c.GetPlanesPeriodo)
	c.Mapping("GetEvaluacion", c.GetEvaluacion)
	c.Mapping("PlanesAEvaluar", c.PlanesAEvaluar)
	c.Mapping("Unidades", c.Unidades)
	c.Mapping("Avances", c.Avances)
}

// GetPlanesPeriodo ...
// @Title GetPlanesPeriodo
// @Description get Planes y vigencias para la unidad y vigencia dado
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	unidad 		path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /planes-periodo/:vigencia/:unidad [get]
func (c *EvaluacionController) GetPlanesPeriodo() {
	defer errorhandler.HandlePanic(&c.Controller)

	vigencia := c.Ctx.Input.Param(":vigencia")
	unidad := c.Ctx.Input.Param(":unidad")

	resultado, err := services.GetPlanesPeriodo(vigencia, unidad)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// GetEvaluacion ...
// @Title GetEvaluacion
// @Description get Evaluacion
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	plan 		path 	string	true		"The key for staticblock"
// @Param	periodo 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /:vigencia/:plan:/:periodo [get]
func (c *EvaluacionController) GetEvaluacion() {
	defer errorhandler.HandlePanic(&c.Controller)

	vigencia := c.Ctx.Input.Param(":vigencia")
	plan := c.Ctx.Input.Param(":plan")
	periodoId := c.Ctx.Input.Param(":periodo")

	resultado, err := services.GetEvaluacion(vigencia, plan, periodoId)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}

// Get Planes A Evaluar ...
// @Title GetPlanesAEvaluar
// @Description get Planes que se pueden evaluar
// @Success 200
// @Failure 404
// @router /planes/ [get]
func (c *EvaluacionController) PlanesAEvaluar() {
	defer errorhandler.HandlePanic(&c.Controller)

	resultado, err := services.PlanesAEvaluar()

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// Get Unidades ...
// @Title GetUnidades
// @Description get Unidades
// @Param	plan 		path 	string	true		"The key for staticblock"
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /unidades/:plan:/:vigencia [get]
func (c *EvaluacionController) Unidades() {
	defer errorhandler.HandlePanic(&c.Controller)

	plan := c.Ctx.Input.Param(":plan")
	vigencia := c.Ctx.Input.Param(":vigencia")

	resultado, err := services.Unidades(plan, vigencia)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}
	c.ServeJSON()
}

// Get Avance ...
// @Title GetAvance
// @Description get Avance de Unidad
// @Param	plan 		path 	string	true		"The key for staticblock"
// @Param	vigencia 	path 	string	true		"The key for staticblock"
// @Param	unidad 	path 	string	true		"The key for staticblock"
// @Success 200
// @Failure 404
// @router /avance/:plan:/:vigencia/:unidad [get]
func (c *EvaluacionController) Avances() {
	defer errorhandler.HandlePanic(&c.Controller)

	plan := c.Ctx.Input.Param(":plan")
	vigencia := c.Ctx.Input.Param(":vigencia")
	unidad := c.Ctx.Input.Param(":unidad")

	resultado, err := services.Avances(plan, vigencia, unidad)

	if err == nil {
		c.Ctx.Output.SetStatus(200)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 200, resultado)
	} else {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = requestresponse.APIResponseDTO(true, 404, nil, err.Error())
	}

	c.ServeJSON()
}
