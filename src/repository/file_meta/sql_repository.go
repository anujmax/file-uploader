package file_meta

import (
	"github.com/anujmax/file-uploader/src/datasources/file_meta"
	"github.com/anujmax/file-uploader/src/domain"
	"log"
)

const (
	queryInsertFileMeta = "INSERT INTO file_metadata(file_identifier, file_name, file_size, file_type, created_date) " +
		"VALUES(?, ?, ?, ?, ? );"
	queryGetFileMeta = "SELECT file_identifier, file_name, file_size, file_type, created_date " +
		"FROM file_metadata WHERE file_identifier=?;"
)

func SaveFileMeta(fileMetadata domain.FileMetaData) error {
	// Create a write transaction
	statement, err := file_meta.Client.Prepare(queryInsertFileMeta)
	if err != nil {
		log.Println("Error creating mysql client" + err.Error())
		return err
	}
	defer statement.Close()

	_, saveErr := statement.Exec(fileMetadata.FileIdentifier, fileMetadata.FileName, fileMetadata.FileSize,
		fileMetadata.FileType, fileMetadata.DateCreated)
	if saveErr != nil {
		log.Println("Error inserting file metadata" + saveErr.Error())
		return saveErr
	}
	return nil
}

func RetrieveFileMeta(fileIdentifier string) (*domain.FileMetaData, error) {
	var fileMetadata domain.FileMetaData

	statement, err := file_meta.Client.Prepare(queryGetFileMeta)
	if err != nil {
		log.Println("Error creating mysql client" + err.Error())
		return &fileMetadata, err
	}
	defer statement.Close()

	result := statement.QueryRow(fileIdentifier)
	if getErr := result.Scan(&fileMetadata.FileIdentifier, &fileMetadata.FileName, &fileMetadata.FileSize,
		&fileMetadata.FileType, &fileMetadata.DateCreated); getErr != nil {
		log.Println("error when trying to get file by id", getErr)
		return &fileMetadata, getErr
	}
	return &fileMetadata, nil
}
