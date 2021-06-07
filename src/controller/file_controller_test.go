package controller

import (
	"encoding/json"
	"errors"
	"github.com/anujmax/file-uploader/src/domain"
	"github.com/anujmax/file-uploader/src/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/png"
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

func TestErrorSavingFile(t *testing.T) {
	service.AuthenticationService = &authenticationServiceMock{}
	service.FileService = &fileServiceMock{}

	authenticate = func(authToken string) error {
		return nil
	}

	saveFile = func(file multipart.File, header multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError) {
		return nil, domain.NewUploadError(
			"Error saving file",
			http.StatusBadRequest,
		)
	}
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("uploadfile", "someimg.png")
		if err != nil {
			t.Error(err)
		}
		img := createImage()

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()
	request := httptest.NewRequest("POST", "/upload", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	r := gin.Default()

	rr := httptest.NewRecorder()
	r.POST("/upload", UploadFile)
	r.ServeHTTP(rr, request)

	var response domain.Response
	_ = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.EqualValues(t, rr.Body.String(), "{\"message\":\"Error saving file\"}")
	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
}

func TestUploadFileSuccess(t *testing.T) {
	service.AuthenticationService = &authenticationServiceMock{}
	service.FileService = &fileServiceMock{}

	authenticate = func(authToken string) error {
		return nil
	}
	fileIdentifier := uuid.New().String()
	saveFile = func(file multipart.File, header multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError) {
		var fileMetadata = &domain.FileMetaData{}
		fileMetadata.FileIdentifier = fileIdentifier
		return fileMetadata, nil
	}
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()
		part, err := writer.CreateFormFile("uploadfile", "someimg.png")
		if err != nil {
			t.Error(err)
		}
		img := createImage()

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}
	}()
	request := httptest.NewRequest("POST", "/upload", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	r := gin.Default()

	rr := httptest.NewRecorder()
	r.POST("/upload", UploadFile)
	r.ServeHTTP(rr, request)

	var response domain.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.EqualValues(t, "Your file has been successfully uploaded.", response.Message)
	assert.EqualValues(t, "/download/"+fileIdentifier, response.Location)
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusCreated, rr.Code)
}

func createImage() *image.RGBA {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}
	// Encode as PNG.
	f, _ := os.Create("image.png")
	png.Encode(f, img)
	return img
}
