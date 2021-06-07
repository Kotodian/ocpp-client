package main

import (
	_ "github.com/joho/godotenv/autoload"
	"ocpp-client/initialize"
	"os"
)

func main() {
	_ = initialize.GinEngine().Run(":" + os.Getenv("PORT"))
}
