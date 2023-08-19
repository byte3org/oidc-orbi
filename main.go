package main

import (
	"github.com/byte3org/oidc-orbi/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	_ = bootstrap.RootApp.Execute()
}
