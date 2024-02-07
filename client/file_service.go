package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	. "shared/interfaces"
	. "shared/types"

	resty "github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type FileService struct {
	RoutePrefix        string `default:"./files"`
	fileProvider       IFileProvider
	fileHashIterator   IFileHashIterator
	merkleTreeProvider IMerkleTreeProvider
	hashProvider       IHashProvider
}

type DownloadFileResponse struct {
	File  []byte
	Proof *Proof
}

type Context struct {
	client *resty.Client
	// That is only for demoing purposes to show that proof works
	isCorruptionNeeded bool
}

func ParseMultipartResponse(resp *resty.Response) (*DownloadFileResponse, error) {
	// I expect here only 2 parts
	// 1. File
	// 2. Proof
	downloadResponse := &DownloadFileResponse{
		Proof: &Proof{},
	}
	fileBuffer := new(bytes.Buffer)

	_, params, err := mime.ParseMediaType(resp.Header().Get("Content-Type"))
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error parsing media type: %v", err))
	}
	bytesReader := bytes.NewReader(resp.Body())
	reader := multipart.NewReader(bytesReader, params["boundary"])

	part, err := reader.NextPart()
	for err == nil {
		switch part.FormName() {
		case "file":
			_, err = io.Copy(fileBuffer, part)
			if err != nil {
				return nil, err
			}
			downloadResponse.File = fileBuffer.Bytes()
			break
		case "proof":
			proofBytes, err := io.ReadAll(part)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(proofBytes, downloadResponse.Proof)
			break
		}
		part, err = reader.NextPart()
	}

	return downloadResponse, nil
}

func (f *FileService) Get(filename string, c *Context) ([]byte, error) {
	log.Debug().Msg(fmt.Sprintf("Downloading file: %s", filename))

	resp, err := c.client.R().
		Get(f.RoutePrefix + "/" + filename)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(resp.String())
	}
	// Parse multipart response
	downloadResponse, err := ParseMultipartResponse(resp)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting stored root hash: %v", err))
	}
	rootHash, err := f.fileHashIterator.GetStoredRootHash()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting stored root hash: %v", err))
	}
	if rootHash != downloadResponse.Proof.RootHash {
		return nil, errors.New(fmt.Sprintf("Root hash mismatch: %s != %s", rootHash, downloadResponse.Proof.RootHash))
	}
	// Such a hack to show that proof works
	if c.isCorruptionNeeded {
		// Corrupt file
		downloadResponse.File = []byte("Corrupted file")
	}

	// Verify proof
	targetHash := f.hashProvider.Hash(downloadResponse.File)
	log.Debug().Msg(fmt.Sprintf("Root hash: %s", downloadResponse.Proof.RootHash))
	log.Debug().Msg(fmt.Sprintf("Proof: %v", downloadResponse.Proof.Hashes))
	log.Debug().Msg(fmt.Sprintf("Target hash: %s", targetHash))
	isValid, err := f.merkleTreeProvider.VerifyProof(targetHash, downloadResponse.Proof)
	if !isValid {
		return nil, errors.New("Invalid proof. Seems like File is corrupted")
	}

	return downloadResponse.File, nil
}

func (f *FileService) Upload(c *Context) error {
	var request = c.client.R()
	// Get list of files inside directory
	// Add files to request
	files, err := f.fileProvider.List()
	for _, file := range files {
		fileBytes, err := f.fileProvider.GetFile(file)
		if err != nil {
			return err
		}
		request.SetFileReader("files", filepath.Base(file), bytes.NewReader(fileBytes))
	}

	// Compute Merkle Tree
	f.merkleTreeProvider.BuildTree()
	log.Info().Msg(fmt.Sprintf("Merkle Root Hash: %s", f.merkleTreeProvider.GetRootHash()))
	log.Debug().Msg(fmt.Sprintf("Uploading files to: %v", f.RoutePrefix))

	// Send request
	resp, err := request.Post(f.RoutePrefix)
	if err != nil {
		return errors.New(fmt.Sprintf("Error uploading files: %v", err))
	}

	if resp.StatusCode() != http.StatusOK {
		log.Info().Msg("Upload response: " + resp.String())
		return errors.New(resp.String())
	}

	// Store root hash
	err = f.fileHashIterator.StoreRootHash(f.merkleTreeProvider.GetRootHash())
	if err != nil {
		return err
	}

	// Remove files
	for _, file := range files {
		err = f.fileProvider.RemoveFile(file)
		if err != nil {
			return err
		}
	}

	log.Info().Msg("Upload response: " + resp.String())
	return nil
}

func NewFileService(
	fileProvider IFileProvider,
	fileHashIterator IFileHashIterator,
	hashProvider IHashProvider,
	merkleTreeProvider IMerkleTreeProvider) *FileService {

	return &FileService{
		fileProvider:       fileProvider,
		fileHashIterator:   fileHashIterator,
		merkleTreeProvider: merkleTreeProvider,
		hashProvider:       hashProvider,
		RoutePrefix:        "/files",
	}
}
