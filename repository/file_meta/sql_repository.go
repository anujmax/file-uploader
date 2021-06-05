package file_meta

import (
	"github.com/anujmax/file-uploader/datasources/memdb"
	"github.com/anujmax/file-uploader/domain"
)

func SaveFileMeta(fileMetadata domain.FileMetaData) (*domain.FileMetaData, error) {
	// Create a write transaction
	txn := memdb.FileDb.Txn(true)
	if err := txn.Insert("file_metadata", fileMetadata); err != nil {
		return nil, err
	}
	// Commit the transaction
	txn.Commit()
	return &fileMetadata, nil
}

func RetrieveFileMeta(FileIdentifier string) (*domain.FileMetaData, error) {
	txn := memdb.FileDb.Txn(false)
	defer txn.Abort()

	// Lookup by identifier
	raw, err := txn.First("file_metadata", "id", FileIdentifier)
	if err != nil {
		return nil, err
	}

	return raw.(*domain.FileMetaData), nil
}
