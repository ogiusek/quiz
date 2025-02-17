package common

import (
	"errors"
)

var (
	ErrRepositoryNotFound             error = errors.New("this model does not exist")
	ErrRepositoryParallelModification error = errors.New("cannot update model because of parrarel modification")
	ErrRepositoryConflict             error = errors.New("conflict when adding model")
)
