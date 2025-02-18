package route

import (
	"gapi/controller/svc1"
	"gapi/controller/svc2"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	var app *gin.Engine = gin.New()

	app.GET("/svc1/req1", svc1.Req1)
	app.GET("/svc1/req2", svc1.Req2)
	app.GET("/svc2/req1", svc2.Req1)
	app.GET("/svc2/req2", svc2.Req2)

	return app
}
