package http

import (
	"fmt"
)

type ClientErrorResponse struct {
	Reason           string
	StatusCode       int
	DeveloperMessage string
	Retry            bool
}

func (e *ClientErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s (dev msg: %s)", e.StatusCode, e.Reason, e.DeveloperMessage)
}
