package database

import (
	"fmt"
	"sync"

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
		})
		if err != nil {
			panic(err)
		}

		// Create Type
		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE participant_career_interest AS ENUM (
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
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE participant_status AS ENUM (
				'waiting-for-verification',
				'verified',
				'declined'
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE membership_role AS ENUM (
				'leader',
				'member'
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE team_category AS ENUM (
				'competitive-programming',
				'datavidia',
				'uxvidia',
				'arkalogica'
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE team_status AS ENUM (
				'waiting-for-evaluation',
				'passed',
				'eliminated'
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE photo_type AS ENUM (
				'pribadi',
				'kartu-pelajar',
				'bukti-mahasiswa-aktif'
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE photo_status AS ENUM (
				'waiting-for-approval',
				'approved',
				'denied'
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

		if err := db.Exec(`
		DO $$ BEGIN
			CREATE TYPE submission_stage AS ENUM (
				'first-stage',
				'second-stage',
				'final-stage'
			);
		EXCEPTION
			WHEN duplicate_object THEN NULL;
		END $$
	`).Error; err != nil {
			panic(err)
		}

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
