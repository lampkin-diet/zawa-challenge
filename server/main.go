package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	. "server/route"
	. "shared/config"
	. "shared/provider"
)

func serve() {
	// Get Config
	config := LoadConfig()
	SetupLogger(config)

	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// Preps
	hashProvider := NewSha256HashProvider()

	fileProvider := NewFileProvider(config.StoragePath)
	// Check that paths are correct
	if fileProvider == nil {
		log.Fatal().Msg("Path " + config.StoragePath + " to directory with files does not exist. Please configure STORAGE_PATH in .env file correctly")
		return
	}
	fileHashIterator := NewFileHashIterator(hashProvider, fileProvider)
	merkleTreeProvider := NewMerkleTreeProvider(fileHashIterator)

	fileRouter := NewFileRouter(fileProvider, hashProvider, merkleTreeProvider)
	
	// File routes
	e.GET("/files/:filename", func(c echo.Context) error {
		return fileRouter.Get(c)
	})
	e.POST("/files", func(c echo.Context) error {
		return fileRouter.Post(c)
	})

	e.Debug = true
	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(fmt.Sprintf("%s:%s", config.Address, config.Port)))
}

func main() {
	serve()
}