package service

import (
	domain2 "github.com/anujmax/file-uploader/src/domain"
	file_meta2 "github.com/anujmax/file-uploader/src/repository/file_meta"
	file_repo2 "github.com/anujmax/file-uploader/src/repository/file_repo"
	utils2 "github.com/anujmax/file-uploader/src/utils"
	"mime/multipart"
	"net/http"
	"time"
)

func SaveFile(file multipart.File, header multipart.FileHeader) (*domain2.FileMetaData, *domain2.UploadError) {
	fileType, err := utils2.GetFileType(file)
	if err != nil {
		return nil, domain2.NewUploadError(
			"Error reading file",
			http.StatusBadRequest,
		)
	}
	if !utils2.IsFileImage(fileType) {
		return nil, domain2.NewUploadError(
			"File is not of type image",
			http.StatusBadRequest,
		)
	}
	fileId, err := file_repo2.Save(file, header)
	if err != nil {
		return nil, domain2.NewUploadError(
			"Unable to save the file",
			http.StatusInternalServerError,
		)
	}
	var fileMetadata = getFileMeta(fileType, fileId, header)
	err = file_meta2.SaveFileMeta(fileMetadata)
	if err != nil {
		return nil, domain2.NewUploadError(
			"Unable to save the file metadata",
			http.StatusInternalServerError,
		)
	}
	return &fileMetadata, nil
}

func getFileMeta(fileType string, fileId string, header multipart.FileHeader) domain2.FileMetaData {
	var fileMetadata domain2.FileMetaData
	fileMetadata.FileType = fileType
	fileMetadata.FileName = header.Filename
	fileMetadata.FileSize = header.Size
	fileMetadata.FileIdentifier = fileId
	fileMetadata.DateCreated = time.Now().UTC().Format("2006-01-02 15:04:05")
	return fileMetadata
}
