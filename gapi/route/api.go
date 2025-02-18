package route

import (
	"gapi/controller/svc1"
	"gapi/controller/svc2"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	var app *gin.Engine = gin.New()

	app_svc1 := app.Group("/svc1")
	app_svc1.GET("/req1", svc1.Req1)
	app_svc1.GET("/req2", svc1.Req2)

	app_svc2 := app.Group("/svc2")
	app_svc2.GET("/req1", svc2.Req1)
	app_svc2.GET("/req2", svc2.Req2)

	return app
}
