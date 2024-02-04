package shared

import "mime/multipart"

type IFileProvider interface {
	FileExists(filename string) (bool, error)
	// GetFile returns the file content
	GetFile(filename string) ([]byte, error)
	// WriteFile writes the file content
	WriteFile(filename string, content []byte) error
	// WriteMultipartFile writes the file content from a multipart file
	WriteMultipartFile(file *multipart.FileHeader) error
	// List returns the list of files in a directory
	List(filename string) ([]string, error)
}
