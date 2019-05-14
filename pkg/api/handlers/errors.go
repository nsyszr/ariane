package handlers

import "fmt"

const (
	ErrInvalidObject int = iota + 1
	/*ErrCodeKeyExists
	ErrCodeResourceVersionConflicts
	ErrCodeInvalidObj
	ErrCodeUnreachable*/
)

var errCodeToMessage = map[int]string{
	ErrInvalidObject: "invalid object",
	/*ErrCodeKeyExists:   "key exists",
	ErrCodeResourceVersionConflicts: "resource version conflicts",
	ErrCodeInvalidObj:               "invalid object",
	ErrCodeUnreachable:              "server unreachable",*/
}

var errCodeToStatusCode = map[int]int{
	ErrInvalidObject: 400,
}

func NewInvalidObjectError(reason string) *HandlerError {
	return &HandlerError{
		Code:    ErrInvalidObject,
		Message: errCodeToMessage[ErrInvalidObject],
		Reason:  reason,
	}
}

type HandlerError struct {
	Code    int
	Message string
	Reason  string
}

// Error method implements the error interface
func (e *HandlerError) Error() string {
	return fmt.Sprintf("StorageError: %s: %s",
		errCodeToMessage[e.Code], e.Reason)
}

func (e *HandlerError) StatusCode() int {
	if _, ok := errCodeToStatusCode[e.Code]; !ok {
		return 500
	}
	return errCodeToStatusCode[e.Code]
}

// IsInvalidObject returns true if and only if err is invalid object error.
func IsInvalidObject(err error) bool {
	return isErrCode(err, ErrInvalidObject)
}

func isErrCode(err error, code int) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*HandlerError); ok {
		return e.Code == code
	}
	return false
}
