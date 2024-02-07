package provider

import (
	"log"
	. "shared/interfaces"
)

type FileHashIterator struct {
	Position int
	HashProvider IHashProvider
	FileProvider IFileProvider
}

func (f *FileHashIterator) Next() (string, bool) {
	files := f.List()
	if f.Position >= len(f.List()) {
		return "", false
	}
	file := files[f.Position]
	// Get file hash
	hash, err := f.GetFileHash(file)
	if err != nil {
		return "", false
	}
	// Update position
	f.Position++
	return hash, true
}

func (f *FileHashIterator) List() []string {
	files, err := f.FileProvider.List()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return files
}

func (f *FileHashIterator) Empty() bool {
	return f.Position >= len(f.List())
}

func (f *FileHashIterator) Reset() {
	f.Position = 0
}

func (f *FileHashIterator) StoreRootHash(rootHash string) error {
	return f.FileProvider.WriteFile("root_hash", []byte(rootHash))
}

func (f *FileHashIterator) GetStoredRootHash() (string, error) {
	b, err := f.FileProvider.GetFile("root_hash")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (f *FileHashIterator) GetListHashes() []string {
	defer f.Reset()

	hashes := make([]string, 0, len(f.List()))
	for hash, ok := f.Next(); ok; hash, ok = f.Next() {
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
	return &FileHashIterator{
		Position: 0,
		HashProvider: hashProvider,
		FileProvider: fileProvider,
	}
}