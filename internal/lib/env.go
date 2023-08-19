package lib

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

type Env struct {
	LogLevel string `mapstructure:"LOG_LEVEL"  validate:"required"`

	DBUsername string `mapstructure:"DB_USER" validate:"required"`
	DBPassword string `mapstructure:"DB_PASS" validate:"required"`
	DBHost     string `mapstructure:"DB_HOST" validate:"required"`
	DBName     string `mapstructure:"DB_NAME" validate:"required"`
	DriverName string `mapstructure:"DRIVER_NAME"`
}

var globalEnv = Env{}

func NewEnv() *Env {
	envVars := make(map[string]string)

	for _, key := range []string{
		"LOG_LEVEL",
		"SERVER_PORT",
		"ENVIRONMENT",
		"DB_USER",
		"DB_PASS",
		"DB_HOST",
		"DB_NAME",
		"DRIVER_NAME",
	} {
		envVars[key] = os.Getenv(key)
	}

	err := mapstructure.Decode(envVars, &globalEnv)

	err = validator.New().Struct(globalEnv)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		log.Fatal(validationErrors)
	}

	return &globalEnv
}
