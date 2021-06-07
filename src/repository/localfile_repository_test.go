package repository

import (
	"github.com/anujmax/file-uploader/src/domain"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"os"
	"testing"
)

func TestSaveFileMeta(t *testing.T) {
	defer os.Remove("test_file.jpg")
	defer os.RemoveAll(FileBasePath)
	file, err := os.Create("test_file.jpg")
	defer file.Close()

	header := multipart.FileHeader{}
	header.Filename = "test_file.jpg"

	fileIdentifier, err := FileRepo.Save(file, header)
	assert.NotEmpty(t, fileIdentifier)
	assert.Nil(t, err)

	fileMeta := domain.FileMetaData{}
	fileMeta.FileName = header.Filename
	fileMeta.FileIdentifier = fileIdentifier
	retrievedFile, err := os.Stat(FileBasePath+"/"+fileIdentifier+"/"+header.Filename)
	assert.NotEmpty(t, retrievedFile)
	assert.Nil(t, err)
	assert.Equal(t, retrievedFile.Name(), header.Filename)
}
