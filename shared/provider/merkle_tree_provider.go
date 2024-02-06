package provider

import (
	. "shared"
	. "shared/merkle"
)

type MerkleTreeProvider struct {
	FileHashIterator IFileHashIterator
	Tree             IMerkleTree
}

func (m *MerkleTreeProvider) GetRootHash() string {
	return m.Tree.GetRootHash()
}

func (m *MerkleTreeProvider) MakeProof(filename string) (*Proof, error) {

	targetHash, err := m.FileHashIterator.GetFileHash(filename)
	if err != nil {
		return nil, err
	}
	return m.Tree.MakeProof(targetHash)
}

func (m *MerkleTreeProvider) VerifyProof(targetHash string, proof *Proof) (bool, error) {
	return m.Tree.VerifyProof(targetHash, proof), nil
}

func (m *MerkleTreeProvider) BuildTree() error {
	m.Tree = NewMerkleTree(m.FileHashIterator.GetListHashes(), m.FileHashIterator.GetHashProvider())
	return nil
}

func NewMerkleTreeProvider(fileHashIterator IFileHashIterator) *MerkleTreeProvider {
	var hashes []string
	var tree = &MerkleTree{
		Root:         nil,
		HashProvider: fileHashIterator.GetHashProvider(),
	}

	hash, ok := fileHashIterator.Next()

	for ok {
		hashes = append(hashes, hash)
		hash, ok = fileHashIterator.Next()
	}

	if len(hashes) > 0 {
		tree = NewMerkleTree(hashes, fileHashIterator.GetHashProvider())
	}

	return &MerkleTreeProvider{
		FileHashIterator: fileHashIterator,
		Tree:             tree,
	}
}
