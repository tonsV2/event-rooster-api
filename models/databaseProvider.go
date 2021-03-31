package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func ProvideDatabase() *gorm.DB {
	database := provideSqliteDatabase()
	//	database := providePostgresDatabase()

	_ = database.AutoMigrate(
		&Event{},
		&Group{},
		&Runner{},
	)

	return database
}

func providePostgresDatabase() *gorm.DB {
	host, _ := os.LookupEnv("DATABASE_HOST")
	port, _ := os.LookupEnv("DATABASE_PORT")
	username, _ := os.LookupEnv("DATABASE_USERNAME")
	password, _ := os.LookupEnv("DATABASE_PASSWORD")
	name, _ := os.LookupEnv("DATABASE_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, username, password, name, port)

	config := gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	database, err := gorm.Open(postgres.Open(dsn), &config)

	if err != nil {
		panic("Failed to connect to database!")
	}

	return database
}

func provideSqliteDatabase() *gorm.DB {
	config := gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	database, err := gorm.Open(sqlite.Open("test.db"), &config)

	if err != nil {
		panic("Failed to connect to database!")
	}

	_ = database.AutoMigrate(
		&Event{},
		&Group{},
		&Runner{},
	)

	return database
}
