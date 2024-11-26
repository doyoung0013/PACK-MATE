package db

import (
	"database/sql"
	"fmt"

	"github.com/GDG-on-Campus-KHU/SC4_BE/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(conf *config.DBConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	// 연결 테스트
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	// 커넥션 풀 설정
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)

	return nil
}
