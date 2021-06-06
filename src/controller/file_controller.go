package controller

import (
	service2 "github.com/anujmax/file-uploader/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadFile(c *gin.Context) {
	authToken := c.Request.FormValue("token")
	authError := service2.Authenticate(authToken)
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
	_, saveError := service2.SaveFile(file, *header)
	if saveError != nil {
		c.JSON(saveError.Status(), gin.H{
			"message": saveError.Message(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Your file has been successfully uploaded.",
	})
}

func DownloadFile(c *gin.Context) {

	c.JSON(http.StatusNotImplemented, gin.H{
		"message": "Not implemented",
	})
}
