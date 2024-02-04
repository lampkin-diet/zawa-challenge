package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	. "shared"

	resty "github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type FileService struct {
	RoutePrefix  string `default:"./files"`
	fileProvider IFileProvider
}

type Context struct {
	client *resty.Client
}

func (f *FileService) Get(filename string, c *Context) ([]byte, error) {
	log.Debug().Msg(fmt.Sprintf("Getting file: %s", filename))

	resp, err := c.client.R().
		Get(f.RoutePrefix + "/" + filename)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(resp.String())
	}
	return resp.Body(), nil
}

func (f *FileService) UploadDir(path string, c *Context) error {
	var request = c.client.R()
	// Get list of files inside directory
	files, err := f.fileProvider.List(path)
	if err != nil {
		return err
	}
	// Add files to request
	for _, file := range files {
		fileBytes, err := f.fileProvider.GetFile(file)
		if err != nil {
			return err
		}
		request.SetFileReader("files", filepath.Base(file), bytes.NewReader(fileBytes))
	}

	log.Debug().Msg(fmt.Sprintf("Uploading files to: %v", f.RoutePrefix))

	// Send request
	resp, err := request.Post(f.RoutePrefix)

	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		log.Info().Msg("Upload response: " + resp.String())
		return errors.New(resp.String())
	}

	log.Debug().Msg("Upload response: " + resp.String())
	return nil
}

func NewFileService(fileProvider IFileProvider) *FileService {
	return &FileService{fileProvider: fileProvider, RoutePrefix: "/files"}
}
