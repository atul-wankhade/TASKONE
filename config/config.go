package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port      string
	JWTSecret string
}

var AppConfig *Config

func LoadConfig() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found , loading environment vsriables.")
	}

	viper.AutomaticEnv()

	AppConfig = &Config{
		Port:      viper.GetString("PORT"),
		JWTSecret: viper.GetString("JWT_SECRET"),
	}

	if AppConfig.Port == "" || AppConfig.JWTSecret == "" {
		log.Fatal("Required environment variables are missing.")
		os.Exit(1)
	}
}
