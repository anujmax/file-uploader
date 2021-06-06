package file_meta

import (
	file_meta2 "github.com/anujmax/file-uploader/src/datasources/file_meta"
	domain2 "github.com/anujmax/file-uploader/src/domain"
	"log"
)

const (
	queryInsertFileMeta = "INSERT INTO file_metadata(file_identifier, file_name, file_size, file_type, created_date) " +
		"VALUES(?, ?, ?, ?, ? );"
	queryGetFileMeta = "SELECT file_identifier, file_name, file_size, file_type, file_size, created_date " +
		"FROM file_metadata WHERE file_identifier=?;"
)

func SaveFileMeta(fileMetadata domain2.FileMetaData) error {
	// Create a write transaction
	statement, err := file_meta2.Client.Prepare(queryInsertFileMeta)
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

func RetrieveFileMeta(FileIdentifier string) (*domain2.FileMetaData, error) {
	/*txn := mysql.GetDb().Txn(false)
	defer txn.Abort()

	// Lookup by identifier
	raw, err := txn.First("file_metadata", "id", FileIdentifier)
	if err != nil {
		return nil, err
	}

	return raw.(*domain.FileMetaData), nil*/
	return nil, nil
}
