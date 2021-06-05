package memdb

import (
	"github.com/hashicorp/go-memdb"
	"log"
)

var (
	FileDb *memdb.MemDB
)

func init() {
	// Create the DB schema
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"file_uploader": &memdb.TableSchema{
				Name: "file_uploader",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "FileIdentifier"},
					},
					"file_name": &memdb.IndexSchema{
						Name:    "file_name",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "FileName"},
					},
					"file_size": &memdb.IndexSchema{
						Name:    "file_size",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "FileSize"},
					},
					"date_created": &memdb.IndexSchema{
						Name:    "date_created",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "DateCreated"},
					},
				},
			},
		},
	}

	FileDb, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
	_ = FileDb
	log.Println("Database successfully created!!")
}
