package domain

type FileMetaData struct {
	FileIdentifier string
	FileName       string
	FileSize       int64
	DateCreated    string
}

type UploadError struct {
	message string
	status  int
}

func (u UploadError) Message() string {
	return u.message
}

func (u UploadError) Status() int {
	return u.status
}

func NewRestError(message string, status int) *UploadError {
	return &UploadError{
		message: message,
		status:  status,
	}
}
