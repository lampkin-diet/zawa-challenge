package route

import (
	"fmt"
	. "shared"
	. "shared/provider"
	. "shared/merkle"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FileRouter struct {
	// FileProvider
	fileProvider IFileProvider
	hashProvider IHashProvider
	fileHashIterator IFileHashIterator
}

func (f *FileRouter) Get(c echo.Context) error {
	filename := c.Param("filename")
	// Check if file exists
	isExist, err := f.fileProvider.FileExists(filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error while checking file existence")
	}
	// If file does not exist, return 404
	if !isExist {
		return c.String(http.StatusNotFound, fmt.Sprintf("File %s not found", filename))
	}
	// Server the file
	return c.File(filename)
}

func (f *FileRouter) Post(c echo.Context) error {
	// Check if file exists
	log.Info().Msg("Uploading files...")
	fileNames := []string{}

	parsed, err := c.MultipartForm()
	if err != nil {
		return c.String(http.StatusBadRequest, "Error while parsing form")
	}
	log.Info().Msg(fmt.Sprintf("Parsed: %v", parsed))
	files := parsed.File["files"]

	log.Info().Msg(fmt.Sprintf("Files: %v", files))

	if err != nil {
		return c.String(http.StatusBadRequest, "No files found")
	}
	for _, file := range files {
		isExist, err := f.fileProvider.FileExists(file.Filename)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Internal Server Error while checking file existence")
		}
		if isExist {
			return c.String(http.StatusBadRequest, fmt.Sprintf("File %s already exists", file.Filename))
		}
		// Write file
		err = f.fileProvider.WriteMultipartFile(file)
		if err != nil {
			log.Error().Err(err).Msg("Error while writing file: " + err.Error())
			return c.String(http.StatusInternalServerError, "Internal Server Error while writing file")
		}
		fileNames = append(fileNames, file.Filename)
	}
	// Generate MerkleTree
	fileHashIterator := NewFileHashIterator(fileNames, f.hashProvider, f.fileProvider)
	merkleTree := NewMerkleTree(fileHashIterator)

	log.Info().Msg("Merkle tree generated: " + merkleTree.Root.Hash)

	return c.String(http.StatusOK, "Files were uploaded successfully")
}

func NewFileRouter(fileProvider IFileProvider) *FileRouter {
	return &FileRouter{
		fileProvider: fileProvider,
		hashProvider: NewSha256HashProvider(),
	}
}
