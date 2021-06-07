package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/anujmax/file-uploader/src/domain"
	"github.com/anujmax/file-uploader/src/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	authenticate func(string) error
	initialize   func(string)
	saveFile     func(multipart.File, multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError)
	retrieveFile func(string) ([]byte, *domain.FileMetaData, *domain.UploadError)
)

type authenticationServiceMock struct{}

func (sm *authenticationServiceMock) Initialize(authToken string) {
	initialize(authToken)
}

func (sm *authenticationServiceMock) Authenticate(authToken string) error {
	return authenticate(authToken)
}

type fileServiceMock struct{}

func (fm *fileServiceMock) SaveFile(file multipart.File, header multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError) {
	return saveFile(file, header)
}

func (fm *fileServiceMock) RetrieveFile(fileIdentifier string) ([]byte, *domain.FileMetaData, *domain.UploadError) {
	return retrieveFile(fileIdentifier)
}

func TestAuthenticationFailure(t *testing.T) {
	service.AuthenticationService = &authenticationServiceMock{}
	service.FileService = &fileServiceMock{}

	authenticate = func(authToken string) error {
		return errors.New("Authentication failure")
	}

	r := gin.Default()
	req, _ := http.NewRequest(http.MethodPost, "/upload", nil)

	rr := httptest.NewRecorder()
	r.POST("/upload", UploadFile)
	r.ServeHTTP(rr, req)

	assert.EqualValues(t, rr.Body.String(), "{\"message\":\"Authentication failure\"}")
	assert.EqualValues(t, http.StatusForbidden, rr.Code)
}

func TestNoFilePresent(t *testing.T) {
	service.AuthenticationService = &authenticationServiceMock{}
	service.FileService = &fileServiceMock{}

	authenticate = func(authToken string) error {
		return nil
	}

	r := gin.Default()
	req, _ := http.NewRequest(http.MethodPost, "/upload", nil)

	rr := httptest.NewRecorder()
	r.POST("/upload", UploadFile)
	r.ServeHTTP(rr, req)

	assert.EqualValues(t, rr.Body.String(), "{\"message\":\"No file is received\"}")
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
}

func TestUploadFile(t *testing.T) {
	service.AuthenticationService = &authenticationServiceMock{}
	service.FileService = &fileServiceMock{}

	authenticate = func(authToken string) error {
		return nil
	}
	fileIdentifier := uuid.New().String()
	saveFile = func(file multipart.File, header multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError) {
		var fileMetadata *domain.FileMetaData
		fileMetadata.FileIdentifier = fileIdentifier
		return fileMetadata, nil
	}
	path := "../../img.png"

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("uploadfile", "test.png")

	sample, err := os.Open(path)
	_, err = io.Copy(part, sample)

	r := gin.Default()
	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	r.POST("/upload", UploadFile)
	r.ServeHTTP(rr, req)

	var response domain.Response
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.EqualValues(t, "Your file has been successfully uploaded.", response.Message)
	assert.EqualValues(t, "/download/"+fileIdentifier, response.Location)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusOK, rr.Code)
}
