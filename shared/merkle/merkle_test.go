package merkle_test

import (
	"fmt"
	. "shared/types"
	. "shared/merkle"
	. "shared/provider"
	"testing"
)

func TestBuildMerkleTree(t *testing.T) {
	hashProvider := &Sha256HashProvider{}
	hashes := []string{"hash1", "hash2", "hash3", "hash4"}
	tree := NewMerkleTree(hashes, hashProvider)
	if tree.Root == nil {
		t.Error("Root is nil")
	}

	if hashProvider.Hash2Nodes(hashes[0], hashes[1]) != tree.Root.Left.Hash {
		t.Error("Left hash is not correct")
	}
	if hashProvider.Hash2Nodes(hashes[2], hashes[3]) != tree.Root.Right.Hash {
		t.Error("Right hash is not correct")
	}

	if hashProvider.Hash2Nodes(tree.Root.Left.Hash, tree.Root.Right.Hash) != tree.Root.Hash {
		t.Error("Root hash is not correct")
	}
}

func TestBuildMerkleTreeOddCount(t *testing.T) {
	hashProvider := &Sha256HashProvider{}
	// It should be duplicated with the latest
	hashes := []string{"hash1", "hash2", "hash3"}
	tree := BuildMerkleTree(hashes, hashProvider)
	if tree.Root == nil {
		t.Error("Root is nil")
	}

	if hashProvider.Hash2Nodes(hashes[0], hashes[1]) != tree.Root.Left.Hash {
		t.Error("Left hash is not correct")
	}
	if hashProvider.Hash2Nodes(hashes[2], hashes[2]) != tree.Root.Right.Hash {
		t.Error("Right hash is not correct")
	}

}

//                           root
//               sha256("${sha256("hash1hash2")${sha256("hash3hash4")}")
//    sha256("hash1hash2")            sha256("hash3hash4")
// hash1                 hash2      hash3              hash4

// Expected proof for hash1: [hash2, sha256("hash3hash4")]
func TestMakeProof(t *testing.T) {
	expectedProof := []string{"hash2", NewSha256HashProvider().Hash2Nodes("hash3", "hash4")}
	expectedIndices := []int64{0, 0}
	targetHash := "hash1"

	hashProvider := &Sha256HashProvider{}
	hashes := []string{"hash1", "hash2", "hash3", "hash4"}
	tree := BuildMerkleTree(hashes, hashProvider)
	if tree.Root == nil {
		t.Error("Root is nil")
	}
	// Walk the tree
	proof, _ := tree.MakeProof(targetHash)
	if proof == nil {
		t.Error("Proof is nil")
	}

	for i, hash := range proof.Hashes {
		if hash != expectedProof[i] {
			t.Error(fmt.Sprintf("Proof is not correct: %s", hash))
		}
		if proof.Indices[i] != expectedIndices[i] {
			t.Error(fmt.Sprintf("Indexes is not correct: %d", proof.Indices[i]))
		}
	}
}

func TestVerifyProof(t *testing.T) {
	targetHash := "hash4"
	hashProvider := &Sha256HashProvider{}
	hashes := []string{"hash1", "hash2", "hash3", "hash4"}
	tree := BuildMerkleTree(hashes, hashProvider)
	expectedProof := &Proof{
		Hashes:   []string{"hash3", NewSha256HashProvider().Hash2Nodes("hash1", "hash2")},
		Indices:  []int64{1, 1},
		RootHash: tree.GetRootHash(),
	}
	if tree.Root == nil {
		t.Error("Root is nil")
	}
	// Walk the tree
	proof, _ := tree.MakeProof(targetHash)
	if proof == nil {
		t.Error("Proof is nil")
	}
	// Verify the proof
	if !tree.VerifyProof(targetHash, expectedProof) {
		t.Error("Proof is not correct")
	}
}
