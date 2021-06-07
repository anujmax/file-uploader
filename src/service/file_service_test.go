package service

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/anujmax/file-uploader/src/domain"
	"github.com/anujmax/file-uploader/src/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"testing"
)

var (
	save         func(multipart.File, multipart.FileHeader) (string, error)
	saveFileMeta func(fileMetadata domain.FileMetaData) error
)

type fileRepoMock struct{}

func (fm *fileRepoMock) Save(file multipart.File, header multipart.FileHeader) (string, error) {
	return save(file, header)
}

func (fm *fileRepoMock) Retrieve(domain.FileMetaData) ([]byte, error) {
	return nil, nil
}

type fileMetaRepoMock struct{}

func (fmr *fileMetaRepoMock) SaveFileMeta(fileMetadata domain.FileMetaData) error {
	return saveFileMeta(fileMetadata)
}

func (fmr *fileMetaRepoMock) RetrieveFileMeta(fileIdentifier string) (*domain.FileMetaData, error) {
	return nil, nil
}
func (fmr *fileMetaRepoMock) Initialize(string, string, string, string) *sql.DB {
	return nil
}

func TestSaveFileSizeNotSupported(t *testing.T) {
	file, _ := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0666)
	header := multipart.FileHeader{}
	header.Size = supportedFileSizeBytes + 10
	res, err := FileService.SaveFile(file, header)
	assert.Nil(t, res)
	assert.Equal(t, err, domain.NewUploadError(
		fmt.Sprintf("File is bigger than supported %d Bytes", supportedFileSizeBytes),
		http.StatusBadRequest,
	))
}

func TestSaveFileContentNotImage(t *testing.T) {
	file, _ := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0666)
	header := multipart.FileHeader{}
	header.Size = supportedFileSizeBytes - 10
	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("content-type", "application/json")
	header.Header = mimeHeader
	res, err := FileService.SaveFile(file, header)
	assert.Nil(t, res)
	assert.Equal(t, err, domain.NewUploadError(
		"File is not of type image",
		http.StatusBadRequest,
	))
}

func TestSaveFileSaveFailedinRepo(t *testing.T) {
	file, _ := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0666)
	header := multipart.FileHeader{}
	header.Size = supportedFileSizeBytes - 10
	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("content-type", "image/jpeg")
	header.Header = mimeHeader

	repository.FileRepo = &fileRepoMock{}
	save = func(file multipart.File, header multipart.FileHeader) (string, error) {
		return "", errors.New("Unable to save the file")
	}
	res, err := FileService.SaveFile(file, header)

	assert.Nil(t, res)
	assert.Equal(t, err, domain.NewUploadError(
		"Unable to save the file",
		http.StatusInternalServerError,
	))
}

func TestSaveFileUnableToSaveMeta(t *testing.T) {
	file, _ := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0666)
	header := multipart.FileHeader{}
	header.Size = supportedFileSizeBytes - 10
	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("content-type", "image/jpeg")
	header.Header = mimeHeader
	fileIdentifier := uuid.New().String()
	repository.FileRepo = &fileRepoMock{}
	save = func(file multipart.File, header multipart.FileHeader) (string, error) {
		return fileIdentifier, nil
	}
	saveFileMeta = func(fileMetadata domain.FileMetaData) error {
		return errors.New("Unable to save the file metadata")
	}

	repository.FileMetaRepo = &fileMetaRepoMock{}
	res, err := FileService.SaveFile(file, header)

	assert.Nil(t, res)
	assert.Equal(t, err, domain.NewUploadError(
		"Unable to save the file metadata",
		http.StatusInternalServerError,
	))
}

func TestSaveFileSuccessful(t *testing.T) {
	file, _ := os.OpenFile("", os.O_WRONLY|os.O_CREATE, 0666)
	header := multipart.FileHeader{}
	header.Size = supportedFileSizeBytes - 10
	header.Filename = "test_file"
	mimeHeader := textproto.MIMEHeader{}
	mimeHeader.Set("content-type", "image/jpeg")
	header.Header = mimeHeader
	fileIdentifier := uuid.New().String()
	repository.FileRepo = &fileRepoMock{}
	save = func(file multipart.File, header multipart.FileHeader) (string, error) {
		return fileIdentifier, nil
	}
	saveFileMeta = func(fileMetadata domain.FileMetaData) error {
		return nil
	}

	repository.FileMetaRepo = &fileMetaRepoMock{}
	res, err := FileService.SaveFile(file, header)
	assert.Nil(t, err)
	assert.Equal(t, res.FileIdentifier, fileIdentifier)
	assert.Equal(t, res.FileSize, int64(supportedFileSizeBytes-10))
	assert.Equal(t, res.FileType, "image/jpeg")
	assert.Equal(t, res.FileName, "test_file")
}
