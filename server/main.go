package main

import (
	"fmt"
	"os"

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

	// Echo instance
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// Preps
	cwd, _ := os.Getwd()
	fileProvider := NewFileProvider(cwd + "/storage")
	fileRouter := NewFileRouter(fileProvider)
	proofRouter := NewProofRouter()
	
	// File routes
	e.GET("/files/:filename", func(c echo.Context) error {
		return fileRouter.Get(c)
	})
	e.POST("/files", func(c echo.Context) error {
		return fileRouter.Post(c)
	})
	// Proof routes
	e.GET("/proof/:file", func(c echo.Context) error {
		return proofRouter.Get(c)
	})

	e.Debug = true
	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(fmt.Sprintf("%s:%s", config.Address, config.Port)))
}

func main() {
	serve()
}