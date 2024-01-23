package main

import (
	"github.com/joho/godotenv"
	"log/slog"
)

func init() {
	// only local
	if err := godotenv.Load(".env"); err != nil {
		//log.Fatal("Error loading .env file")
		slog.Warn("Config file not found")
	}
}
