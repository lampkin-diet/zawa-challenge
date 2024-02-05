package provider

import (
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
	// Get file content
	content, err := f.FileProvider.GetFile(file)
	if err != nil {
		return "", false
	}
	// Hash file content
	hash := f.HashProvider.Hash(content)
	// Update position
	f.Position++
	return hash, true
}

func (f *FileHashIterator) Empty() bool {
	return f.Position >= len(f.Files)
}

func (f *FileHashIterator) GetList() []string {
	return f.Files
}

func (f *FileHashIterator) GetFileProvider() IFileProvider {
	return f.FileProvider
}

func (f *FileHashIterator) GetHashProvider() IHashProvider {
	return f.HashProvider
}


func NewFileHashIterator(files []string, hashProvider IHashProvider, fileProvider IFileProvider) *FileHashIterator {
	return &FileHashIterator{
		Files: files,
		Position: 0,
		HashProvider: hashProvider,
		FileProvider: fileProvider,
	}
}