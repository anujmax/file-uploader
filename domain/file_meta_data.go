package domain

type FileMetaData struct {
	FileIdentifier string
	FileName       string
	FileSize       int64
	DateCreated    string
}

type UploadError interface {
	Message() string
	Status() int
}

type uploadError struct {
	message string
	status  int
}

func (u uploadError) Message() string {
	return u.message
}

func (u uploadError) Status() int {
	return u.status
}

func NewRestError(message string, status int) UploadError {
	return uploadError{
		message: message,
		status:  status,
	}
}
