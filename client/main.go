package main

import (
	"fmt"
	"path/filepath"
	. "shared/config"
	. "shared/provider"
	. "shared"

	resty "github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "client",
		Short: "Client for reciving the file with Merkle Proof",
		Long:  "Such CLI helps to upload a lot of files and get the proof for each file respectively",
	}
)

func checkEnv() {
	
}

func main() {

	var (
		fileProvider IFileProvider
		fileService  IFileService
		proofManager IProofManager
	)
	
	config := LoadConfig()
	SetupLogger(config)

	// Create file provider and services
	fileProvider = NewFileProvider(config.StoragePath)
	fileService = NewFileService(fileProvider)
	proofManager = NewProofManager()

	// Check that paths are correct
	pathToFiles, err := filepath.Abs(config.StoragePath)
	log.Debug().Msg(fmt.Sprintf("Path to files: %s", pathToFiles))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get absolute path to storage")
		return
	}
	isExist, err := fileProvider.FileExists(pathToFiles)
	if !isExist {
		log.Fatal().Err(err).Msg("Path to directory with files does not exist. Please configure STORAGE_PATH in .env file correctly")
		return
	}

	context := &Context{
		client: resty.New(),
	}
	log.Debug().Msg(fmt.Sprintf("Server address: %s:%s", config.Address, config.Port))
	// Set base URL
	context.client.SetBaseURL(fmt.Sprintf("http://%s:%s", config.Address, config.Port))

	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload files",
		Long:  "Upload files to the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fileService.UploadDir(pathToFiles, context))
		},
	}

	proofCmd := &cobra.Command{
		Use:   "proof <name of file to proof>",
		Args:  cobra.ExactArgs(1),
		Short: "Get proof for files",
		Long:  "Get proof for files from the server",
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			proof := proofManager.GetProof(filename, context)
			fmt.Println("Proof that file is not changed: ", proof)
		},
	}

	downloadCmd := &cobra.Command{
		Use:   "download <name of file to download>",
		Args:  cobra.ExactArgs(1),
		Short: "Download files",
		Long:  "Download files from the server",
		Run: func(cmd *cobra.Command, args []string) {
			filename := args[0]
			file, err := fileService.Get(filename, context)
			if err != nil {
				fmt.Println(err)
				return
			}
			fileProvider.WriteFile(filename, file)
		},
	}

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(proofCmd)
	rootCmd.AddCommand(downloadCmd)

	rootCmd.Execute()
}
