package main

import (
	"encoding/json"
	controller "library/Controller"
	"library/config"
	"library/database"
	"library/service"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	loadConfig()
	db := database.Connect()
	service.Connect(db)

}
func main() {
	godotenv.Load()
	controller.Routes().Run("localhost:" + config.AppConfig.Port)
}

func loadConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config.AppConfig); err != nil {
		return err
	}

	return nil
}
