package merkle

import (
	. "shared"
)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
}

func (m *MerkleTree) GetRootHash() string {
	return m.Root.Hash
}

func BuildMerkleTree(hashes []string, hashProvider IHashProvider) *MerkleTree {
	var nodes []MerkleNode
	for _, hash := range hashes {
		node := MerkleNode{Hash: hash}
		nodes = append(nodes, node)
	}
	for len(nodes) > 1 {
		var newLevel []MerkleNode
		if len(nodes)%2 != 0 {
			nodes = append(nodes, nodes[len(nodes)-1])
		}
		for i := 0; i < len(nodes); i += 2 {
			node := MerkleNode{
				Left:  &nodes[i],
				Right: &nodes[i+1],
				Hash:  hashProvider.Hash2Nodes(nodes[i].Hash, nodes[i+1].Hash),
			}
			newLevel = append(newLevel, node)
		}
		nodes = newLevel
	}
	return &MerkleTree{Root: &nodes[0]}
}

func NewMerkleTree(fileHashIterator IFileHashIterator) *MerkleTree {
	var hashes []string

	hash, ok := fileHashIterator.Next()

	for ok {
		hash, ok = fileHashIterator.Next()
		hashes = append(hashes, hash)
	}
	if len(hashes) == 0 {
		return nil
	}

	if len(hashes) != len(fileHashIterator.GetList()) {
		return nil
	}

	return BuildMerkleTree(hashes, fileHashIterator.GetHashProvider())
}
