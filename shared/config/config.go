package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port        string `default:"8080"`
	Address     string `default:"localhost"`
	StoragePath string `default:"storage"`
	LogLevel    string `default:"debug"`
}

func LoadConfig() Config {

	port := "8080"
	address := "localhost"
	storagePath := "storage"
	logLevel := "debug"

	err := godotenv.Load()
	if err != nil {
		log.Info().Msgf("Error while reading .env file. environment or default values will be used")
	}

	if os.Getenv("SERVER_PORT") != "" {
		port = os.Getenv("SERVER_PORT")
	}
	if os.Getenv("SERVER_ADDRESS") != "" {
		address = os.Getenv("SERVER_ADDRESS")
	}

	if os.Getenv("STORAGE_PATH") != "" {
		storagePath = os.Getenv("STORAGE_PATH")
	}

	if os.Getenv("LOG_LEVEL") != "" {
		logLevel = os.Getenv("LOG_LEVEL")
	}

	return Config{
		Port:        port,
		Address:     address,
		StoragePath: storagePath,
		LogLevel:    logLevel,
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
