package filesmodule

import "errors"

var errNotAnImageExtension error = errors.New("this is not an image extension")

var imageExtensions []string = []string{
	"",
}

type ImageId FileId // VARCHAR 36 + 1 + 5 = 42

func NewImageId(fileid FileId) ImageId {
	return ImageId(fileid)
}

func (image ImageId) File() FileId {
	return FileId(image)
}

func (imageId *ImageId) Valid() error {
	ext := ((*FileId)(imageId)).FileExtension()
	for _, allowedExt := range imageExtensions {
		if ext == allowedExt {
			return nil
		}
	}
	return errNotAnImageExtension
}

var (
	// TODO
	DefaultUserImage ImageId = ImageId("")
)
