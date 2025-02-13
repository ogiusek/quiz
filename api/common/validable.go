package common

import "fmt"

type ErrorWithPath interface {
	error

	// returns only message
	Message() string

	// returns only path
	Path() string

	// adds property name to path
	// parse property name to which you want to display
	//
	// note: this method modifies and returns modified version for convenience
	Property(string) ErrorWithPath
}

type errorWithPath struct {
	message string
	path    string
}

func (err *errorWithPath) Error() string {
	return fmt.Sprintf("`%s` -> %s", err.path, err.message)
}

func (err *errorWithPath) Message() string {
	return err.message
}

func (err *errorWithPath) Path() string {
	return err.path
}

func (err *errorWithPath) Property(parent string) ErrorWithPath {
	if err.path == "" {
		err.path = parent
	} else {
		err.path = fmt.Sprintf("%s.%s", parent, err.path)
	}
	return err
}

func NewErrorWithPath(message string) ErrorWithPath {
	return &errorWithPath{
		message: message,
		path:    "",
	}
}
func ErrPath(err error) ErrorWithPath {
	if errWithPath, ok := err.(ErrorWithPath); ok {
		return errWithPath
	}
	return &errorWithPath{
		message: err.Error(),
		path:    "",
	}
}

type Validable interface {
	Valid() []error
}
