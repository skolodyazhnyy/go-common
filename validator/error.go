package validator

import (
	"fmt"
	"strings"
)

// Error provides information about validation error
type Error struct {
	// Path to the field related to the error
	Path []string

	// Message describing error
	Message string
}

// Error implements error interface
func (e Error) Error() string {
	return fmt.Sprintf("[%v] %v", strings.Join(e.Path, "."), e.Message)
}
