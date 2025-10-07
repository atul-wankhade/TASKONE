package db

import (
	"TASKONE/config"
	"database/sql"
	"fmt"
	"log"
)

var MySQL *sql.DB

func InitMySQL() {
	MySQL, err := sql.Open("mysql", config.AppConfig.MySQLDSN+"?parseTime=true")
	if err != nil {
		log.Fatal("Failed to connect to MySQL", err)
	}

	if err := MySQL.Ping(); err != nil {
		log.Fatal("MySQL is not reachable", err)
	}
	log.Println("Connected to MySQL")
}

func ConnectMySQL() *sql.DB {
	cfg := config.AppConfig

	// Build DSN safely (escape password if needed)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.MySQLUser,
		cfg.MySQLPassword,
		cfg.MySQLHost,
		cfg.MySQLPort,
		cfg.MySQLDB)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	return db
}
