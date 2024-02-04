package main


type IFileService interface {
	UploadDir(path string, c *Context) error
	Get(filename string, c *Context) ([]byte, error)
}

type IProofManager interface {
	GetProof(filename string, c *Context) (string)
}
