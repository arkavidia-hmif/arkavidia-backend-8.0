package database

import (
	"os"
	"strconv"
	"sync"
)

type DatabaseMetadata struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     int
}

type DatabaseConfig struct {
	metadata DatabaseMetadata
	once     sync.Once
}

// Private
func (databaseConfig *DatabaseConfig) lazyInit() {
	databaseConfig.once.Do(func() {
		host := os.Getenv("POSTGRES_HOST")
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		dbname := os.Getenv("POSTGRES_DBNAME")
		port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
		if err != nil {
			panic(err)
		}

		databaseConfig.metadata.Host = host
		databaseConfig.metadata.User = user
		databaseConfig.metadata.Password = password
		databaseConfig.metadata.DBName = dbname
		databaseConfig.metadata.Port = port
	})
}

// Public
func (databaseConfig *DatabaseConfig) GetMetadata() DatabaseMetadata {
	databaseConfig.lazyInit()
	return databaseConfig.metadata
}

var Config = &DatabaseConfig{}
