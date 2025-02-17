package main

import "github.com/gin-gonic/gin"

func main() {
	app := gin.Default()
	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})
	app.Run("0.0.0.0:8000")
}
