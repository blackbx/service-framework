package response

import (
	"fmt"
	"net/http"
)

type ProblemError interface {
	error
	Problem() *Problem
}

// NewHTTPProblem creates a new instance of a Problem for HTTP errors
func NewHTTPProblem(statusCode int, detail string) *Problem {
	return &Problem{
		Status: statusCode,
		Type:   "https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
		Title:  http.StatusText(statusCode),
		Detail: detail,
	}
}

// Problem is a struct that provides standard error details
type Problem struct {
	Status   int    `json:"status"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Instance string `json:"instance,omitempty"`
}

// Error implements the error interface, which allows a
// problem to be used as an error
func (p *Problem) Error() string {
	return fmt.Sprintf(p.Detail)
}

// Problem implements the ProblemError interface, so that it is
// possible to determine via type assertion whether an error can
// give you a *response.Problem or not
func (p *Problem) Problem() *Problem {
	return p
}
