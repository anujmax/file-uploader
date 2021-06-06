package main

import (
	"github.com/anujmax/file-uploader/controller"
	"github.com/anujmax/file-uploader/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	authToken, err := service.GetAuthToken()
	if err != nil {
		panic(err)
	}
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"auth_token": authToken,
		})
	})
	router.POST("/upload", controller.UploadFile)
	router.POST("/download/:id", controller.DownloadFile)
	router.Run(":8080")
}
