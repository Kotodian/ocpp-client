package main

import (
	_ "github.com/joho/godotenv/autoload"
	"ocpp-client/init"
	"os"
)

func main() {
	_ = init.GinEngine().Run(":" + os.Getenv("PORT"))

}
