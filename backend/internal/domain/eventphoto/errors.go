package eventphoto

import "errors"

var (
	ErrPhotoNotFound      = errors.New("photo not found")
	ErrMaxPhotosReached   = errors.New("maximum of 10 photos per event reached")
	ErrInvalidContentType = errors.New("only image files are allowed (jpeg, png, gif, webp)")
	ErrFileTooLarge       = errors.New("file size exceeds the 10 MB limit")
)
