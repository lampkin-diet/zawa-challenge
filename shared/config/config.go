package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port string
	Address string
	StoragePath string
	LogLevel string
}

func LoadConfig() Config {
	err := godotenv.Load() 
	if err != nil {
		log.Fatal().Msgf("Error loading .env file: %v", err)
		panic(err)
	}

	return Config{
		Port: os.Getenv("SERVER_PORT"),
		Address: os.Getenv("SERVER_ADDRESS"),
		StoragePath: os.Getenv("STORAGE_PATH"),
		LogLevel: os.Getenv("LOG_LEVEL"),
	}
}

func SetupLogger(config Config) {
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	log.Debug().Msgf("Log level changed to : %s", config.LogLevel)
}