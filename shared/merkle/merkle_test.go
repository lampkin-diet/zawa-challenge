package merkle_test

import (
	"testing"
	. "shared/merkle"
	. "shared/provider"
)

func TestBuildMerkleTree(t *testing.T) {
	hashProvider := &Sha256HashProvider{}
	hashes := []string{"hash1", "hash2", "hash3", "hash4"}
	tree := BuildMerkleTree(hashes, hashProvider)
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