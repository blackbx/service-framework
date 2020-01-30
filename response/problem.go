package response

import "net/http"

// NewHTTPProblem creates a new instance of a Problem for HTTP errors
func NewHTTPProblem(statusCode int, detail string) Problem {
	return Problem{
		Status: statusCode,
		Type:   "https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html",
		Title:  http.StatusText(statusCode),
		Detail: detail,
	}
}

// Problem is a struct that provides standard error details
type Problem struct {
	Status int    `json:"status"`
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}
