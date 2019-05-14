package storage

import "fmt"

const (
	ErrCodeKeyNotFound int = iota + 1
	ErrCodeKeyExists
	/*ErrCodeResourceVersionConflicts
	ErrCodeInvalidObj
	ErrCodeUnreachable*/
)

var errCodeToMessage = map[int]string{
	ErrCodeKeyNotFound: "key not found",
	ErrCodeKeyExists:   "key exists",
	/*ErrCodeResourceVersionConflicts: "resource version conflicts",
	ErrCodeInvalidObj:               "invalid object",
	ErrCodeUnreachable:              "server unreachable",*/
}

func NewKeyNotFoundError(key string) *StorageError {
	return &StorageError{
		Code: ErrCodeKeyNotFound,
		Key:  key,
	}
}

func NewKeyExistsError(key string) *StorageError {
	return &StorageError{
		Code: ErrCodeKeyExists,
		Key:  key,
	}
}

type StorageError struct {
	Code int
	Key  string
}

// IsNotFound returns true if and only if err is "key" not found error.
func IsNotFound(err error) bool {
	return isErrCode(err, ErrCodeKeyNotFound)
}

// IsNodeExist returns true if and only if err is an node already exist error.
func IsNodeExist(err error) bool {
	return isErrCode(err, ErrCodeKeyExists)
}

func (e *StorageError) Error() string {
	return fmt.Sprintf("StorageError: %s, Code: %d, Key: %s",
		errCodeToMessage[e.Code], e.Code, e.Key)
}

func isErrCode(err error, code int) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*StorageError); ok {
		return e.Code == code
	}
	return false
}
