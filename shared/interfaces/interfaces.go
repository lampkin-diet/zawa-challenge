package interfaces

import (
	"mime/multipart"
	. "shared/types"
)

type IFileProvider interface {
	FileExists(filename string) (bool, error)
	// GetFile returns the file content
	GetFile(filename string) ([]byte, error)
	// WriteFile writes the file content
	WriteFile(filename string, content []byte) error
	// WriteMultipartFile writes the file content from a multipart file
	WriteMultipartFile(file *multipart.FileHeader) error
	// List returns the list of files in a directory
	List() ([]string, error)
	// RemoveFile removes a file
	RemoveFile(filename string) error
}

type IHashProvider interface {
	Hash(data []byte) string
	Hash2Nodes(left, right string) string
}

type IFileHashIterator interface {
	Next() (string, bool)
	Empty() bool
	Reset()
	StoreRootHash(rootHash string) error
	GetStoredRootHash() (string, error)
	GetFileHash(filename string) (string, error)
	GetListHashes() []string
	GetFileProvider() IFileProvider
	GetHashProvider() IHashProvider
}

type IMerkleTreeProvider interface {
	GetRootHash() string
	MakeProof(filename string) (*Proof, error)
	VerifyProof(targetHash string, proof *Proof) (bool, error)
	BuildTree() error
}

type IMerkleTree interface {
	GetRootHash() string
	MakeProof(targetHash string) (*Proof, error)
	VerifyProof(targetHash string, proof *Proof) bool
}
