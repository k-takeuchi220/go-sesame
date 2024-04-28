package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	deviceUUID := os.Getenv("DEVICE_UUID")

	if err := sendCommand(apiKey, deviceUUID, secretKey, CommandToggle); err != nil {
		fmt.Println(err)
		return
	}

}
