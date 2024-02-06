package provider

import (
	"log"
	. "shared"
)

type FileHashIterator struct {
	Files []string
	Position int
	HashProvider IHashProvider
	FileProvider IFileProvider
}

func (f *FileHashIterator) Next() (string, bool) {
	if f.Position >= len(f.Files) {
		return "", false
	}
	file := f.Files[f.Position]
	// Get file hash
	hash, err := f.GetFileHash(file)
	if err != nil {
		return "", false
	}
	// Update position
	f.Position++
	return hash, true
}

func (f *FileHashIterator) Empty() bool {
	return f.Position >= len(f.Files)
}

func (f *FileHashIterator) Reset() {
	f.Position = 0
}

func (f *FileHashIterator) StoreRootHash(rootHash string) error {
	return f.FileProvider.WriteFile("root_hash", []byte(rootHash))
}

func (f *FileHashIterator) GetRootHash() (string, error) {
	b, err := f.FileProvider.GetFile("root_hash")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *FileHashIterator) GetListHashes() []string {
	hashes := make([]string, 0, len(f.Files))
	for _, file := range f.Files {
		hash, err := f.GetFileHash(file)
		if err != nil {
			log.Fatal(err)
			return nil
		}
		hashes = append(hashes, hash)
	}
	return hashes
}

func (f *FileHashIterator) GetFileProvider() IFileProvider {
	return f.FileProvider
}

func (f *FileHashIterator) GetHashProvider() IHashProvider {
	return f.HashProvider
}

func (f *FileHashIterator) GetFileHash(filename string) (string, error) {
	content, err := f.FileProvider.GetFile(filename)
	if err != nil {
		return "", err
	}
	return f.HashProvider.Hash(content), nil
}


func NewFileHashIterator(hashProvider IHashProvider, fileProvider IFileProvider) *FileHashIterator {
	files, err := fileProvider.List()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &FileHashIterator{
		Files: files,
		Position: 0,
		HashProvider: hashProvider,
		FileProvider: fileProvider,
	}
}