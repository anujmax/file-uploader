package file_repo

import (
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
)

func Save(file multipart.File, header multipart.FileHeader) (string, error) {
	var fileIdentifier = uuid.New().String()
	err := os.Mkdir("./savepath/"+fileIdentifier, 0755)
	if err != nil {
		log.Println("Error creating file path" + err.Error())
		return "", err
	}
	f, err := os.OpenFile("./savepath/"+fileIdentifier+"/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
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

func Retrieve(fileId string) {

}
