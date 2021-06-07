package service

import (
	"fmt"
	"github.com/anujmax/file-uploader/src/domain"
	"github.com/anujmax/file-uploader/src/repository"
	"github.com/anujmax/file-uploader/src/utils"
	"mime/multipart"
	"net/http"
	"time"
)

var (
	FileService fileServiceInterface = &fileService{}
)

type fileService struct{}

type fileServiceInterface interface {
	SaveFile(multipart.File, multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError)
	RetrieveFile(string) ([]byte, *domain.FileMetaData, *domain.UploadError)
}

const supportedFileSizeBytes = 8388608

func (f *fileService) SaveFile(file multipart.File, header multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError) {
	if header.Size > supportedFileSizeBytes {
		return nil, domain.NewUploadError(
			fmt.Sprintf("File is bigger than supported %d Bytes", supportedFileSizeBytes),
			http.StatusBadRequest,
		)
	}

	fileType := header.Header.Get("content-type")
	if len(fileType) < 0 {
		return nil, domain.NewUploadError(
			"Error reading file content",
			http.StatusBadRequest,
		)
	}
	if !utils.IsFileImage(fileType) {
		return nil, domain.NewUploadError(
			"File is not of type image",
			http.StatusBadRequest,
		)
	}
	fileId, err := repository.FileRepo.Save(file, header)
	if err != nil {
		return nil, domain.NewUploadError(
			"Unable to save the file",
			http.StatusInternalServerError,
		)
	}
	var fileMetadata = getFileMeta(fileType, fileId, header)
	err = repository.FileMetaRepo.SaveFileMeta(fileMetadata)
	if err != nil {
		return nil, domain.NewUploadError(
			"Unable to save the file metadata",
			http.StatusInternalServerError,
		)
	}
	return &fileMetadata, nil
}

func getFileMeta(fileType string, fileId string, header multipart.FileHeader) domain.FileMetaData {
	var fileMetadata domain.FileMetaData
	fileMetadata.FileType = fileType
	fileMetadata.FileName = header.Filename
	fileMetadata.FileSize = header.Size
	fileMetadata.FileIdentifier = fileId
	fileMetadata.DateCreated = time.Now().UTC().Format("2006-01-02 15:04:05")
	return fileMetadata
}

func (f *fileService) RetrieveFile(fileIdentifier string) ([]byte, *domain.FileMetaData, *domain.UploadError) {
	fileMeta, err := repository.FileMetaRepo.RetrieveFileMeta(fileIdentifier)
	if err != nil {
		return nil, nil, domain.NewUploadError(
			"Unable to get the file metadata",
			http.StatusNotFound,
		)
	}
	data, err := repository.FileRepo.Retrieve(*fileMeta)
	if err != nil {
		return nil, nil, domain.NewUploadError(
			"Unable to donwload file",
			http.StatusNotFound,
		)
	}
	return data, fileMeta, nil
}
