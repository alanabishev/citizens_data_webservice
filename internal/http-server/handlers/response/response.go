// Package response provides a structure for HTTP responses and functions to create them.
package response

// Response is a structure for HTTP responses.
// It includes a status and an optional error message.
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// StatusOK and StatusError are constants for the status of the response.
const (
	StatusOK    = "OK"
	StatusError = "Error"
)

// OK is a function that creates a successful response.
// It returns a Response with the status set to "OK".
func OK() Response {
	return Response{
		Status: StatusOK,
	}
}

// Error is a function that creates an error response.
// It takes an error message as a parameter and returns a Response with the status set to "Error" and the error message.
func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
