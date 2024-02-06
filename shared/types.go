package shared

type ClientContext struct{}

type Proof struct {
	Hashes   []string `json:"ProofHashes"`
	RootHash string   `json:"RootHash"`
	Indices  []int64  `json:"Indices"`
}
