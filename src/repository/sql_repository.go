package repository

import (
	"database/sql"
	"fmt"
	"github.com/anujmax/file-uploader/src/domain"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	FileMetaRepo fileMetaRepoInterface = &fileMetaRepo{}
)

const (
	queryInsertFileMeta = "INSERT INTO file_metadata(file_identifier, file_name, file_size, file_type, created_date) " +
		"VALUES(?, ?, ?, ?, ? );"
	queryGetFileMeta = "SELECT file_identifier, file_name, file_size, file_type, created_date " +
		"FROM file_metadata WHERE file_identifier=?;"
)

/**
Repository to help Save and retrieve file meta like size, type
*/
type fileMetaRepoInterface interface {
	SaveFileMeta(fileMetadata domain.FileMetaData) error
	RetrieveFileMeta(fileIdentifier string) (*domain.FileMetaData, error)
	Initialize(string, string, string, string) *sql.DB
}

type fileMetaRepo struct {
	db *sql.DB
}

func (fm *fileMetaRepo) Initialize(username, password, host, schema string) *sql.DB {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username, password, host, schema,
	)
	log.Println("datasource= ", dataSourceName)
	var err error
	fm.db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	if err = fm.db.Ping(); err != nil {
		panic(err)
	}
	log.Println("database successfully configured")
	return fm.db
}

func (fm *fileMetaRepo) SaveFileMeta(fileMetadata domain.FileMetaData) error {
	// Create a write transaction
	statement, err := fm.db.Prepare(queryInsertFileMeta)
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

func (fm *fileMetaRepo) RetrieveFileMeta(fileIdentifier string) (*domain.FileMetaData, error) {
	var fileMetadata domain.FileMetaData

	statement, err := fm.db.Prepare(queryGetFileMeta)
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
