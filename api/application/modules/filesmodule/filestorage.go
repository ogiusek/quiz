package filesmodule

import "errors"

var (
	ErrFileAlreadyExists  error = errors.New("file with this id already exists")
	ErrFileDoNotExist     error = errors.New("file with given id do not exists")
	ErrCanceledFileUpload error = errors.New("file upload were canceled")
)

type FileStorage interface {
	GetFile(FileId) (File, error)
	UploadFile(File, func(fileId FileId) bool) (FileId, error)
	Track(FileId) error
	UnTrack(FileId) error
}

// TODO
