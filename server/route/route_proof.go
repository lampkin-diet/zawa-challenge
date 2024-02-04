package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ProofRouter struct {
}

func (f *ProofRouter) Get(c echo.Context) error {
	return c.String(http.StatusOK, "ProofRouter Get")
}


func NewProofRouter() *ProofRouter {
	return &ProofRouter{}
}
