package main


type IFileService interface {
	Upload(c *Context) error
	Get(filename string, c *Context) ([]byte, error)
}
