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
            Method: "Avances",
            Router: "/avance/:plan:/:vigencia/:unidad",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "GetPlanesPeriodo",
            Router: "/planes-periodo/:vigencia/:unidad",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "PlanesAEvaluar",
            Router: "/planes/",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"] = append(beego.GlobalControllerRouter["github.com/udistrital/planeacion_evaluacion_mid/controllers:EvaluacionController"],
        beego.ControllerComments{
            Method: "Unidades",
            Router: "/unidades/:plan:/:vigencia",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
