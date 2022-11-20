package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	databaseConfig "arkavidia-backend-8.0/competition/config/database"
	"arkavidia-backend-8.0/competition/models"
)

var currentDB *gorm.DB = nil

func Init() *gorm.DB {
	config := databaseConfig.GetDatabaseConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", config.Host, config.User, config.Password, config.DBName, config.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	// Create Type
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DROP TYPE IF EXISTS participant_career_interest`).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			`CREATE TYPE participant_career_interest AS ENUM (
				'software-engineering',
				'product-management',
				'ui-designer',
				'ux-designer',
				'ux-researcher',
				'it-consultant',
				'game-developer',
				'cyber-security',
				'business-analyst',
				'business-intelligence',
				'data-scientist',
				'data-analyst'
			)`,
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DROP TYPE IF EXISTS team_category`).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			`CREATE TYPE team_category AS ENUM (
				'competitive-programming',
				'datavidia',
				'uxvidia',
				'arkalogica'
			)`,
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DROP TYPE IF EXISTS membership_role`).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			`CREATE TYPE membership_role AS ENUM (
				'leader',
				'member-1',
				'member-2'
			)`,
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DROP TYPE IF EXISTS photo_status`).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			`CREATE TYPE photo_status AS ENUM (
				'waiting-for-verification',
				'verified',
				'declined'
			)`,
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DROP TYPE IF EXISTS photo_type`).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			`CREATE TYPE photo_type AS ENUM (
				'pribadi',
				'kartu-pelajar',
				'bukti-mahasiswa-aktif'
			)`,
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`DROP TYPE IF EXISTS submission_stage`).Error; err != nil {
			return err
		}
		if err := tx.Exec(
			`CREATE TYPE submission_stage AS ENUM (
				'first-stage',
				'second-stage',
				'final-stage'
			)`,
		).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}

	// Migrate Class
	if err := db.AutoMigrate(&models.Participant{}, &models.Team{}, &models.Membership{}, &models.Photo{}, &models.Submission{}); err != nil {
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
