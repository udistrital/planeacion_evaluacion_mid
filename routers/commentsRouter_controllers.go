package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "GetEvaluacion",
            Router: "/:vigencia/:plan:/:periodo",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "GetPlanesPeriodo",
            Router: "/planes_periodo/:vigencia/:unidad",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:MainController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:MainController"],
        beego.ControllerComments{
            Method: "Get",
            Router: "/get/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
