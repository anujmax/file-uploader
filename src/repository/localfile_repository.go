package repository

import (
	"fmt"
	"github.com/anujmax/file-uploader/src/domain"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
)

var (
	FileRepo fileRepoInterface = &fileRepo{}
)

type fileRepo struct{}

/**
Repository to help Save and retrieve file contents
 */
type fileRepoInterface interface {
	Save(multipart.File, multipart.FileHeader) (string, error)
	Retrieve(domain.FileMetaData) ([]byte, error)
}

const FileBasePath = "./savepath"

func (fr *fileRepo) Save(file multipart.File, header multipart.FileHeader) (string, error) {
	var fileIdentifier = uuid.New().String()
	_, err := os.Stat(FileBasePath)
	if err != nil {
		err := os.Mkdir(FileBasePath, 0755)
		if err != nil {
			log.Println("Error creating file path" + err.Error())
			return "", err
		}
	}
	err = os.Mkdir(FileBasePath+"/"+fileIdentifier, 0755)
	if err != nil {
		log.Println("Error creating file path" + err.Error())
		return "", err
	}

	f, err := os.OpenFile(FileBasePath+"/"+fileIdentifier+"/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println("Error closing file pointer" + err.Error())
		return "", err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("Error closing file pointer" + err.Error())
		}
	}(f)
	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}
	return fileIdentifier, nil
}

func (fr *fileRepo) Retrieve(fileMeta domain.FileMetaData) ([]byte, error) {
	data, err := ioutil.ReadFile(FileBasePath + "/" + fileMeta.FileIdentifier + "/" + fileMeta.FileName)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}
	return data, nil
}
