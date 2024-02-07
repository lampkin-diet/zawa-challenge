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
	StorageDir string
}

func (f *FileProvider) GetFullPath(filename string) string {
	return filepath.Join(f.StorageDir, filename)
}

func (f *FileProvider) GetFile(filename string) ([]byte, error) {
	log.Debug().Msg(fmt.Sprintf("Getting file: %s", filename))
	fpath := f.GetFullPath(filename)
	return os.ReadFile(fpath)
}

func (f *FileProvider) WriteFile(filename string, file []byte) error {
	path := f.GetFullPath(filename)
	return os.WriteFile(path, file, os.ModePerm)
}

func (f *FileProvider) FileExists(filename string) (bool, error) {
	fpath := f.GetFullPath(filename)
	log.Debug().Msg(fmt.Sprintf("Checking file/dir: %s", fpath))
	stat, err := os.Stat(fpath)
	if err == nil {
		return true, nil
	}
	log.Debug().Msg(fmt.Sprintf("Error: %v", err))
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

func (f *FileProvider) List() ([]string, error) {
	files, err := os.ReadDir(f.StorageDir)
	if err != nil {
		return nil, err
	}
	var list []string
	for _, file := range files {
		list = append(list, file.Name())
	}
	return list, nil
}

func (f *FileProvider) RemoveFile(filename string) error {
	return os.Remove(f.GetFullPath(filename))
}

func NewFileProvider(storagePrefix string) *FileProvider {
	fpath, err := filepath.Abs(storagePrefix)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while getting absolute path: " + err.Error())
		return nil
	}
	return &FileProvider{StorageDir: fpath}
}