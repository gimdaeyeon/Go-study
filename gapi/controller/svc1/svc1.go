package svc1

import (
	"gapi/model"

	"github.com/gin-gonic/gin"
)

func Req1(c *gin.Context) {
	result := model.GetAdminList()
	// c.String(200, result)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.String(200, result)

}

func Req2(c *gin.Context) {
	c.JSON(200, gin.H{
		"SVC1": "REQ2",
	})
}
