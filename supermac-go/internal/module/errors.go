package module

import "fmt"

const (
	ExitOK         = 0
	ExitGeneral    = 1
	ExitUsage      = 2 // Bad command/flag usage
	ExitPermission = 3 // sudo/permission denied
	ExitNetwork    = 4 // Network unreachable
	ExitNotFound   = 5 // Resource not found
)

// ExitError represents a structured error with an exit code.
type ExitError struct {
	Code    int
	Message string
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// NewExitError creates a structured exit error.
func NewExitError(code int, msg string) *ExitError {
	return &ExitError{Code: code, Message: msg}
}
