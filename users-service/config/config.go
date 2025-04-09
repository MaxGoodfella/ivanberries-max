package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DBURI           string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func LoadConfig() *Config {
	dbHost := os.Getenv("POSTGRES_HOST")
	if dbHost == "" {
		log.Fatal("POSTGRES_HOST is required")
	}

	dbUser := os.Getenv("POSTGRES_USER")
	if dbUser == "" {
		log.Fatal("POSTGRES_USER is required")
	}

	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	if dbPassword == "" {
		log.Fatal("POSTGRES_PASSWORD is required")
	}

	dbName := os.Getenv("POSTGRES_DB")
	if dbName == "" {
		log.Fatal("POSTGRES_DB is required")
	}

	dbPort := os.Getenv("POSTGRES_PORT")
	if dbPort == "" {
		log.Fatal("POSTGRES_PORT is required")
	}

	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	accessTTL := getEnvAsInt("ACCESS_TOKEN_TTL", 15)
	refreshTTL := getEnvAsInt("REFRESH_TOKEN_TTL", 10080)

	return &Config{
		DBURI:           dbURI,
		JWTSecret:       jwtSecret,
		AccessTokenTTL:  time.Duration(accessTTL) * time.Minute,
		RefreshTokenTTL: time.Duration(refreshTTL) * time.Minute,
	}
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Fatalf("Invalid value for %s: %v", name, err)
	}
	return val
}
