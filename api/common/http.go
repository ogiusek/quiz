package common

type HttpError interface {
	error
	StatusCode() int
}

type httpError struct {
	message    string
	statusCode int
}

func (err httpError) Error() string {
	return err.message
}

func (err httpError) StatusCode() int {
	return err.statusCode
}

func NewHttpError(message string, statusCode int) HttpError {
	return httpError{
		message:    message,
		statusCode: statusCode,
	}
}

// errors

var (
	ErrHttpParallelModification error = NewHttpError("parallel modification try again", 409)
	ErrHttpServerError          error = NewHttpError("server error", 500)
)
