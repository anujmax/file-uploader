package service

import (
	"github.com/anujmax/file-uploader/domain"
	"github.com/anujmax/file-uploader/repository/file_meta"
	"github.com/anujmax/file-uploader/repository/file_repo"
	"github.com/anujmax/file-uploader/utils"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"time"
)

func SaveFile(c *gin.Context) (*domain.FileMetaData, *domain.UploadError) {
	file, header, err := c.Request.FormFile("uploadfile")
	if err != nil {
		return nil, domain.NewUploadError(
			"No file is received",
			http.StatusBadRequest,
		)
	}
	fileType, err := utils.GetFileType(file)
	if err != nil {
		return nil, domain.NewUploadError(
			"Error reading file",
			http.StatusBadRequest,
		)
	}
	if !utils.IsFileImage(fileType) {
		return nil, domain.NewUploadError(
			"File is not of type image",
			http.StatusBadRequest,
		)
	}

	fileId, err := file_repo.Save(file, *header)
	if err != nil {
		return nil, domain.NewUploadError(
			"Unable to save the file",
			http.StatusInternalServerError,
		)
	}
	var fileMetadata = getFileMeta(fileType, fileId, *header)
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
