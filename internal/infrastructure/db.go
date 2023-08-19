package infrastructure

import (
	"fmt"
	"log"

	"github.com/byte3org/oidc-orbi/internal/lib"
	"github.com/byte3org/oidc-orbi/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database modal
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(env *lib.Env) Database {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s  sslmode=%s TimeZone=UTC",
		env.DBHost, env.DBUsername, env.DBPassword, env.DBName, "disable")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		log.Panic(err)
	}

	err = db.AutoMigrate(&models.User{})

	if err != nil {
		log.Panic(err)
	}

	return Database{DB: db}
}
