package file_repo

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

const FileBasePath = "./savepath"

func Save(file multipart.File, header multipart.FileHeader) (string, error) {
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
	//tempFile, err := ioutil.TempFile(FileBasePath+"/"+fileIdentifier, header.Filename)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//defer tempFile.Close()
	//fileBytes, err := ioutil.ReadAll(file)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//tempFile.Write(fileBytes)

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

func Retrieve(fileMeta domain.FileMetaData) ([]byte, error) {
	data, err := ioutil.ReadFile(FileBasePath + "/" + fileMeta.FileIdentifier + "/" + fileMeta.FileName)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, err
	}
	return data, nil
}
