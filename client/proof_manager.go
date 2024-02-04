package main

type ProofManager struct {
	RoutePrefix string `default:"/proof"`
}

func (f *ProofManager) GetProof(filename string, c *Context) (string) {
	_, err := c.client.R().
		Get(f.RoutePrefix + filename)
	if err != nil {
		return ""
	}
	return ""
}

func NewProofManager() *ProofManager {
	return &ProofManager{}
}