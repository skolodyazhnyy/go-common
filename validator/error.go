package validator

import (
	"fmt"
	"strings"
)

// Collection of errors
type Errors []error

// Error implements error interface
func (e Errors) Error() string {
	var msgs []string

	for _, err := range e {
		msgs = append(msgs, err.Error())
	}

	return strings.Join(msgs, "\n")
}

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
