package app

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DebugLevel int
	DbURL      string
}

func NewConfig() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port := ":" + os.Getenv("PORT")
	debugLevel, err := strconv.Atoi(os.Getenv("DEBUG_LEVEL"))
	if err != nil {
		log.Fatalf("Can`t parse config field debugLevel: %s", err)
	}
	dbURL := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_SSL_MODE"))
	fmt.Println("Database URL: ", dbURL)
	return &Config{port, debugLevel, dbURL}
}
