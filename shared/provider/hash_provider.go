package provider

import (
	"crypto/sha256"
	"encoding/hex"
)

type Sha256HashProvider struct{}

func (h *Sha256HashProvider) Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (h *Sha256HashProvider) Hash2Nodes(left, right string) string {
	hash := sha256.Sum256([]byte(left + right))
	return hex.EncodeToString(hash[:])
}

func NewSha256HashProvider() *Sha256HashProvider {
	return &Sha256HashProvider{}
}