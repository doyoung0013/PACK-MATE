package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func GetDBConfig() *DBConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"), //두 번째 인자는 기본값
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "gdgdb"),
	}
}

// 환경변수를 가져오는 헬퍼 함수
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
