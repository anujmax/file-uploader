package controller

import (
	"github.com/anujmax/file-uploader/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UploadFile(c *gin.Context) {
	authToken := c.Request.FormValue("token")
	authError := service.Authenticate(authToken)
	if authError != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"message": authError.Error(),
		})
		return
	}
	_, saveError := service.SaveFile(c)
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
