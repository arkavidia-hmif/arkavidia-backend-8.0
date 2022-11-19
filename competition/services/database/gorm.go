package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	databaseConfig "arkavidia-backend-8.0/competition/config/database"
	"arkavidia-backend-8.0/competition/models"
)

var currentDB *gorm.DB = nil

func Init() *gorm.DB {
	databaseConfig := databaseConfig.GetDatabaseConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", databaseConfig.Host, databaseConfig.User, databaseConfig.Password, databaseConfig.DBName, databaseConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.Participant{}, &models.Team{}, &models.Membership{}, &models.Photo{}, &models.Submission{})
	if err != nil {
		panic(err)
	}

	return db
}

func GetDB() *gorm.DB {
	if currentDB == nil {
		currentDB = Init()
	}

	return currentDB
}
