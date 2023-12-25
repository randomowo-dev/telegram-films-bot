package client

import (
	"fmt"
)

type ErrorResponse struct {
	Reason           string
	StatusCode       int
	DeveloperMessage string
	Retry            bool
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%d: %s (dev msg: %s)", e.StatusCode, e.Reason, e.DeveloperMessage)
}
