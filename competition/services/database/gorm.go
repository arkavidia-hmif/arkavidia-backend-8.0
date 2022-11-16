package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"arkavidia-backend-8.0/competition/models"
)

var currentDB *gorm.DB = nil

func Init() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DBNAME"), os.Getenv("POSTGRES_PORT"))
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
