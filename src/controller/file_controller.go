package controller

import (
	"fmt"
	"github.com/anujmax/file-uploader/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadFile(c *gin.Context) {
	authToken := c.Request.FormValue("token")
	authError := service.AuthenticationService.Authenticate(authToken)
	if authError != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": authError.Error(),
		})
		return
	}
	file, header, err := c.Request.FormFile("uploadfile")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
	}
	fileMetaData, saveError := service.FileService.SaveFile(file, *header)
	if saveError != nil {
		c.JSON(saveError.Status(), gin.H{
			"message": saveError.Message(),
		})
		return
	}
	location := "/download/" + fileMetaData.FileIdentifier
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Your file has been successfully uploaded.",
		"Location": location,
	})
	c.Header("Location", "/download/"+fileMetaData.FileIdentifier)
}

func DownloadFile(c *gin.Context) {
	fileIdentifier := c.Param("id")
	data, fileMeta, err := service.FileService.RetrieveFile(fileIdentifier)
	if err != nil {
		c.JSON(err.Status(), gin.H{
			"message": err.Message(),
		})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileMeta.FileName))
	c.Data(http.StatusOK, fileMeta.FileType, data)

}
