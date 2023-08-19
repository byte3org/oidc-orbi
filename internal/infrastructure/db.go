package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database modal
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(logger lib.Logger, env *lib.Env) Database {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s  sslmode=%s TimeZone=UTC",
		env.DBHost, env.DBUsername, env.DBPassword, env.DBName, "disable")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		logger.Panic(err)
	}

	/*
		err = db.AutoMigrate(&models.Blog{},
			&models.Answer{}, &models.Forum{},
			&models.User{}, &models.Answer{},
			&models.Like{},
			&models.Event{}, &models.Speaker{},
		)

		if err != nil {
			logger.Panic(err)
		}
	*/

	return Database{DB: db}
}
