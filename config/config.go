package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port          string `mapstructure:"PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	MySQLDSN      string `mapstructure:"MYSQL_DSN"`
	MongoURI      string `mapstructure:"MONGO_URI"`
	MongoDB       string `mapstructure:"MONGO_DB"`
	MySQLUser     string `mapstructure:"MYDQL_USER"`
	MySQLPassword string `mapstructure:"MYSQL_PASSWORD"`
	MySQLHost     string `mapstructure:"MYSQL_HOST"`
	MySQLPort     string `mapstructure:"MYSQL_PORT"`
	MySQLDB       string `mapstructure:"MYSQL_DB"`
}

var AppConfig *Config

func LoadConfig() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found , loading environment vsriables.")
	}

	viper.AutomaticEnv()

	AppConfig = &Config{
		Port:          viper.GetString("PORT"),
		JWTSecret:     viper.GetString("JWT_SECRET"),
		MySQLDSN:      viper.GetString("MYSQL_DSN"),
		MongoURI:      viper.GetString("MONGO_URI"),
		MongoDB:       viper.GetString("MONGO_DB"),
		MySQLUser:     viper.GetString("MYSQL_USER"),
		MySQLPassword: viper.GetString("MYSQL_PASSWORD"),
		MySQLHost:     viper.GetString("MYSQL_HOST"),
		MySQLPort:     viper.GetString("MYSQL_PORT"),
		MySQLDB:       viper.GetString("MYSQL_DB"),
	}

	log.Printf("Loaded Config: %+v", AppConfig)

	if AppConfig.Port == "" || AppConfig.JWTSecret == "" {
		log.Fatal("Required environment variables are missing.")
		os.Exit(1)
	}

	if AppConfig.MongoDB == "" || AppConfig.MongoURI == "" {
		log.Fatal("Missing DB Config.")
		os.Exit(1)
	}
}
