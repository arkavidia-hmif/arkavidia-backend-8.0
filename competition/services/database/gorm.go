package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	databaseConfig "arkavidia-backend-8.0/competition/config/database"
	"arkavidia-backend-8.0/competition/models"
)

type Database struct {
	connection *gorm.DB
	once       sync.Once
}

// Private
func (database *Database) lazyInit() {
	database.once.Do(func() {
		config := databaseConfig.Config.GetMetadata()
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", config.Host, config.User, config.Password, config.DBName, config.Port)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				return time.Now().Local()
			},
			PrepareStmt: true,
		})
		if err != nil {
			panic(err)
		}

		// Register Plugins
		db.Use(Plugins)

		// Migrate Class
		if err := db.AutoMigrate(&models.Participant{}, &models.Team{}, &models.Membership{}, &models.Photo{}, &models.Submission{}); err != nil {
			panic(err)
		}

		// Assign To Struct
		database.connection = db
	})
}

// Public
func (database *Database) GetConnection() *gorm.DB {
	database.lazyInit()
	return database.connection
}

var DB = &Database{}
