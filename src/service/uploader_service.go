package service

import (
	"fmt"
	"github.com/anujmax/file-uploader/src/domain"
	"github.com/anujmax/file-uploader/src/repository/file_meta"
	"github.com/anujmax/file-uploader/src/repository/file_repo"
	"github.com/anujmax/file-uploader/src/utils"
	"mime/multipart"
	"net/http"
	"time"
)

const supportedFileSizeBytes = 8388608

func SaveFile(file multipart.File, header multipart.FileHeader) (*domain.FileMetaData, *domain.UploadError) {
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
	fileId, err := file_repo.Save(file, header)
	if err != nil {
		return nil, domain.NewUploadError(
			"Unable to save the file",
			http.StatusInternalServerError,
		)
	}
	var fileMetadata = getFileMeta(fileType, fileId, header)
	err = file_meta.SaveFileMeta(fileMetadata)
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

func RetrieveFile(fileIdentifier string) ([]byte, *domain.FileMetaData, *domain.UploadError) {
	fileMeta, err := file_meta.RetrieveFileMeta(fileIdentifier)
	if err != nil {
		return nil, nil, domain.NewUploadError(
			"Unable to get the file metadata",
			http.StatusNotFound,
		)
	}
	data, err := file_repo.Retrieve(*fileMeta)
	if err != nil {
		return nil, nil, domain.NewUploadError(
			"Unable to donwload file",
			http.StatusNotFound,
		)
	}
	return data, fileMeta, nil
}
