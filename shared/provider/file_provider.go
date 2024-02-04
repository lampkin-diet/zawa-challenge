package provider

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type FileProvider struct {
	StoragePrefix string
}

func (f *FileProvider) GetFile(filename string) ([]byte, error) {
	fpath := filepath.Join(f.StoragePrefix, filename)
	return os.ReadFile(fpath)
}

func (f *FileProvider) WriteFile(filename string, file []byte) error {
	path := filepath.Join(f.StoragePrefix, filename)
	return os.WriteFile(path, file, os.ModePerm)
}

func (f *FileProvider) FileExists(fpath string) (bool, error) {
	log.Debug().Msg(fmt.Sprintf("Checking file/dir: %s", fpath))
	stat, err := os.Stat(fpath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	if stat.IsDir() {
		return true, nil
	}
	return false, err
}

func (f *FileProvider) WriteMultipartFile(file *multipart.FileHeader) error {
	log.Debug().Msg(fmt.Sprintf("Writing file: %s", file.Filename))
	// Open file
	sourceFile, err := file.Open()
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	// Read file
	fileContent, err := io.ReadAll(sourceFile)
	if err != nil {
		return err
	}
	// Write file
	return f.WriteFile(file.Filename, fileContent)
}

func (f *FileProvider) List(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, file := range files {
		list = append(list, file.Name())
	}
	return list, nil
}

func NewFileProvider(storagePrefix string) *FileProvider {
	return &FileProvider{StoragePrefix: storagePrefix}
}