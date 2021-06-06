package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func GetFileType(file multipart.File) (string, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}
	filetype := http.DetectContentType(buf.Bytes())
	return filetype, nil
}

func IsFileImage(fileType string) bool {
	switch fileType {
	case "image/jpeg", "image/jpg":
		return true
	case "image/gif":
		return true
	case "image/png":
		return true
	default:
		return false
	}
}
