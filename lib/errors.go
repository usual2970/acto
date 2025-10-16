package lib

import "fmt"

// ErrorType classifies library errors.
type ErrorType int

const (
	ErrTypeConnection ErrorType = iota
	ErrTypeValidation
	ErrTypeConfiguration
	ErrTypeRepository
	ErrTypeService
	ErrTypeRoute
)

func (t ErrorType) String() string {
	switch t {
	case ErrTypeConnection:
		return "connection"
	case ErrTypeValidation:
		return "validation"
	case ErrTypeConfiguration:
		return "configuration"
	case ErrTypeRepository:
		return "repository"
	case ErrTypeService:
		return "service"
	case ErrTypeRoute:
		return "route"
	default:
		return "unknown"
	}
}

// LibraryError is a typed error with optional cause and context.
type LibraryError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *LibraryError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause=%v)", e.Type.String(), e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type.String(), e.Message)
}

func NewLibraryError(t ErrorType, msg string, cause error) *LibraryError {
	return &LibraryError{Type: t, Message: msg, Cause: cause}
}
