package svc1

import "github.com/gin-gonic/gin"

func Req1(c *gin.Context) {
	c.JSON(200, gin.H{
		"SVC1": "REQ1",
	})
}

func Req2(c *gin.Context) {
	c.JSON(200, gin.H{
		"SVC1": "REQ2",
	})
}
