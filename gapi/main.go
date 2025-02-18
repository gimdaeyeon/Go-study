package main

import (
	"gapi/route"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// .env 파일 로드
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env file: %v", err)
	}

	// 환경 변수에서 포트 가져오기
	port := os.Getenv("PORT")

	// go에서 string은 nil이 될 수 없음 따라서 빈 문자열인지만 확인
	// 추가적으로 확인해야한다면 포트번호에 올 수 있는 숫자인지 확인하면 될 것 같다.(정규표현식 활용)
	if port == "" {
		log.Fatalf("PORT is not set in the .env file")
	}

	// app := gin.Default()
	app := route.Router()

	app.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	app.Run("0.0.0.0:" + port)
}
