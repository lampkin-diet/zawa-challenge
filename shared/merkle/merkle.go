package merkle

import (
	"fmt"
	. "shared/types"
	. "shared/interfaces"

	"github.com/rs/zerolog/log"
)

type MerkleTree struct {
	Root         *MerkleNode
	HashProvider IHashProvider
}

type MerkleNode struct {
	Hash  string
	Left  *MerkleNode
	Right *MerkleNode
}

func (m *MerkleTree) GetRootHash() string {
	return m.Root.Hash
}

func (m *MerkleTree) MakeProof(targetHash string) (*Proof, error) {
	proof := []string{}
	// 0 - left, 1 - right
	indices := []int64{}

	// Walk the tree
	var walk func(node *MerkleNode) bool
	walk = func(node *MerkleNode) bool {
		if node.Left == nil && node.Right == nil {
			if node.Hash == targetHash {
				return true
			}
			return false
		}
		if node.Left != nil {
			if walk(node.Left) {
				proof = append(proof, node.Right.Hash)
				indices = append(indices, 0)
				return true
			}
		}
		if node.Right != nil {
			if walk(node.Right) {
				proof = append(proof, node.Left.Hash)
				indices = append(indices, 1)
				return true
			}
		}
		return false
	}

	if !walk(m.Root) {
		return nil, fmt.Errorf("Hash " + targetHash + " not found in the tree")
	}

	return &Proof{Hashes: proof, RootHash: m.GetRootHash(), Indices	: indices}, nil
}

func (m *MerkleTree) VerifyProof(targetHash string, proof *Proof) bool {
	// Index 0 - left, 1 - right
	log.Debug().Msg(fmt.Sprintf("Proof indeces are: %v", proof.Indices))
	hash := targetHash
	for i, index := range proof.Indices {
		if index == 0 {
			hash = m.HashProvider.Hash2Nodes(hash, proof.Hashes[i])
		} else {
			hash = m.HashProvider.Hash2Nodes(proof.Hashes[i], hash)
		}
	}
	return hash == proof.RootHash

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
	return &MerkleTree{Root: &nodes[0], HashProvider: hashProvider}
}

func NewMerkleTree(hashes []string, hashProvider IHashProvider) *MerkleTree {
	log.Debug().Msg(fmt.Sprintf("Building merkle tree with hashes: %v", hashes))
	return BuildMerkleTree(hashes, hashProvider)
}
