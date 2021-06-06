package service

import (
	"github.com/anujmax/file-uploader/domain"
	"github.com/anujmax/file-uploader/repository/file_meta"
	"github.com/anujmax/file-uploader/repository/file_repo"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SaveFile(c *gin.Context) (*domain.FileMetaData, *domain.UploadError) {
	file, header, err := c.Request.FormFile("uploadfile")
	if err != nil {
		return nil, domain.NewRestError(
			"No file is received",
			http.StatusBadRequest,
		)
	}
	var fileMetadata domain.FileMetaData
	fileMetadata.FileName = header.Filename
	fileMetadata.FileSize = header.Size
	fileId, err := file_repo.Save(file, *header)
	if err != nil {
		return nil, domain.NewRestError(
			"Unable to save the file",
			http.StatusInternalServerError,
		)
	}
	fileMetadata.FileIdentifier = fileId
	fileMetadata.DateCreated = time.Now().UTC().Format("2006-01-02 15:04:05")
	err = file_meta.SaveFileMeta(fileMetadata)
	if err != nil {
		return nil, domain.NewRestError(
			"Unable to save the file metadata",
			http.StatusInternalServerError,
		)
	}
	return &fileMetadata, nil
}
