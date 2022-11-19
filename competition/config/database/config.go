package database

import (
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

var currentDatabaseConfig *DatabaseConfig = nil

func Init() *DatabaseConfig {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DBNAME")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		panic(err)
	}

	return &DatabaseConfig{
		Host:     host,
		User:     user,
		Password: password,
		DBName:   dbname,
		Port:     port,
	}
}

func GetDatabaseConfig() *DatabaseConfig {
	if currentDatabaseConfig == nil {
		currentDatabaseConfig = Init()
	}

	return currentDatabaseConfig
}
