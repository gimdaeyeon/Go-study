package model

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DBConn *sql.DB

func Init() {
	var err error

	// .env 파일 로드

	if err = godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 환경 변수에서 DB정보 가져오기
	HOST := os.Getenv("TEST_DB_CONFIG_HOST")
	PORT := os.Getenv("TEST_DB_CONFIG_PORT")
	DBNAME := os.Getenv("TEST_DB_CONFIG_DBNAME")
	USERNAME := os.Getenv("TEST_DB_CONFIG_USERNAME")
	PASSWORD := os.Getenv("TEST_DB_CONFIG_PASSWORD")
	MAX_IDLE_CONNS := os.Getenv("TEST_DB_CONFIG_MAX_IDLE_CONNS")
	MAX_OPEN_CONNS := os.Getenv("TEST_DB_CONFIG_MAX_OPEN_CONNS")

	DSN := USERNAME + ":" + PASSWORD + "@tcp(" + HOST + ":" + PORT + ")/" + DBNAME
	DBConn, err = sql.Open("mysql", DSN)

	if err != nil {
		log.Fatalf("Error conntect db: %v", err)
	}

	// Connection Pool
	_MAX_IDLE_CONNS, _ := strconv.Atoi(MAX_IDLE_CONNS)
	_MAX_OPEN_CONNS, _ := strconv.Atoi(MAX_OPEN_CONNS)
	DBConn.SetMaxIdleConns(_MAX_IDLE_CONNS)
	DBConn.SetMaxOpenConns(_MAX_OPEN_CONNS)
	DBConn.SetConnMaxLifetime(time.Hour)

}

func GetAdminList() string {
	result_str := "["
	var login_id string
	var passwd string
	var nick string
	var email string

	rows, err := DBConn.Query("SELECT LOGIN_ID, PASSWD, NICk, EMAIL FROM TB_ADMIN")
	// rows, err := DBConn.Query("CALL SP_L_ADMIN()")
	if err != nil {
		log.Fatalf("%v", err)
	}
	if rows != nil {
		defer rows.Close()
	}

	cnt := 0
	for rows.Next() {
		if cnt == 0 {
			result_str += "{"
		} else {
			result_str += ",{"
		}
		err := rows.Scan(&login_id, &passwd, &nick, &email)
		if err != nil {
			log.Fatalf("%v", err)
		}
		result_str += `"lOGIN_ID": "` + login_id + `", "PASSWD":"` + passwd + `", "NICK": "` + nick + `", "EMAIL": "` + email + `"}`
		cnt += 1
	}

	result_str += "]"
	return result_str
}
