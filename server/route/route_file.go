package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	
	. "shared/interfaces"
	. "shared/provider"

	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

type FileRouter struct {
	// FileProvider
	fileProvider       IFileProvider
	hashProvider       IHashProvider
	fileHashIterator   IFileHashIterator
	merkleTreeProvider IMerkleTreeProvider
}

func (f *FileRouter) Get(c echo.Context) error {
	filename := c.Param("filename")
	// Check if file exists
	isExist, err := f.fileProvider.FileExists(filename)
	log.Info().Msg(fmt.Sprintf("File %s exists: %v", filename, isExist))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error while checking file existence")
	}
	// If file does not exist, return 404
	if !isExist {
		return c.String(http.StatusNotFound, fmt.Sprintf("File %s not found", filename))
	}
	// Serve the file
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error while creating form file: "+err.Error())
	}
	file, err := f.fileProvider.GetFile(filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error while getting file: "+err.Error())
	}
	_, err = io.Copy(part, bytes.NewReader(file))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error while copying file: "+err.Error())
	}

	// Compile proof for sending
	proof, err := f.merkleTreeProvider.MakeProof(filename)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error while making proof: "+err.Error())
	}
	jsonProof, err := json.Marshal(proof)
	writer.CreateFormField("proof")
	writer.WriteField("proof", string(jsonProof))

	writer.Close()
	c.Response().Header().Set(echo.HeaderContentType, writer.FormDataContentType())
	c.Response().WriteHeader(http.StatusOK)

	_, err = c.Response().Write(body.Bytes())
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal Server Error while writing response")
	}
	return nil
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
	f.merkleTreeProvider.BuildTree()

	log.Info().Msg("Merkle tree generated: " + f.merkleTreeProvider.GetRootHash())

	return c.String(http.StatusOK, "Files were uploaded successfully")
}

func NewFileRouter(fileProvider IFileProvider, hashProvider IHashProvider, merkleTreeProvider IMerkleTreeProvider) *FileRouter {
	// Build tree if there are files
	files, _ := fileProvider.List()
	if len(files) > 0 {
		merkleTreeProvider.BuildTree()
	}
	return &FileRouter{
		fileProvider:       fileProvider,
		hashProvider:       NewSha256HashProvider(),
		merkleTreeProvider: merkleTreeProvider,
	}
}
