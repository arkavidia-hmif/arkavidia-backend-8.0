package database

// TODO: Tambahkan gorm plugin menggunakan go-migration untuk menginisialisasikan type pada PostgreSQL
// REFERENCE: https://gorm.io/docs/write_plugins.html
// REFERENCE: https://github.com/golang-migrate/migrate#use-in-your-go-project
// ASSIGNED TO: @rayhankinan
// STATUS: DONE

import (
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	databaseConfig "arkavidia-backend-8.0/competition/config/database"
)

type MigrationPlugins struct {
	once sync.Once
}

// Private
func (migratorPlugins *MigrationPlugins) lazyInit(connection *gorm.DB) {
	migratorPlugins.once.Do(func() {
		sqlDB, err := connection.DB()
		if err != nil {
			panic(err)
		}

		driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			panic(err)
		}

		config := databaseConfig.Config.GetMetadata()
		m, err := migrate.NewWithDatabaseInstance(
			"file://migrations",
			config.DBName,
			driver,
		)
		if err != nil {
			panic(err)
		}

		m.Up()
	})
}

// Public
func (migrationPlugins MigrationPlugins) Name() string {
	return "migration-plugin"
}

func (migrationPlugins *MigrationPlugins) Initialize(connection *gorm.DB) error {
	migrationPlugins.lazyInit(connection)
	return nil
}

var Plugins = &MigrationPlugins{}
