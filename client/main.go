package main

import (
	"fmt"
	. "shared"
	. "shared/config"
	. "shared/provider"
	"strconv"

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
	)
	
	config := LoadConfig()
	SetupLogger(config)

	// Create file provider and services
	fileProvider = NewFileProvider(config.StoragePath)
	// Check that paths are correct
	if fileProvider == nil {
		log.Fatal().Msg("Path " + config.StoragePath + " to directory with files does not exist. Please configure STORAGE_PATH in .env file correctly")
		return
	}

	hashProvider := NewSha256HashProvider()
	fileHashIterator := NewFileHashIterator(hashProvider, fileProvider)
	merkleTreeProvider := NewMerkleTreeProvider(fileHashIterator)

	fileService = NewFileService(fileProvider, fileHashIterator, hashProvider, merkleTreeProvider)

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
			err := fileService.Upload(context)
			if err != nil {
				fmt.Println(err)
				return
			}
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

	generateFileCmd := &cobra.Command{
		Use:   "generate <number of files to generate>",
		Args: func (c *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("Invalid number of arguments")
			}
			if _, err := strconv.Atoi(args[0]); err != nil {
				return fmt.Errorf("Invalid number of files")
			}
			return nil
		},
		Short: "Generate files",
		Long:  "Generate files to the storage",
		Run: func(cmd *cobra.Command, args []string) {
			count, _ := strconv.Atoi(args[0])
			for i := 0; i < count; i++ {
				fileProvider.WriteFile(fmt.Sprintf("file%d", i), []byte(fmt.Sprintf("file%d", i)))
			}
		},
	}

	rootCmd.AddCommand(uploadCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(generateFileCmd)

	rootCmd.Execute()
}
