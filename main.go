package main

import (
	controller2 "github.com/anujmax/file-uploader/src/controller"
	service2 "github.com/anujmax/file-uploader/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	authToken, err := service2.GetAuthToken()
	if err != nil {
		panic(err)
	}
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"auth_token": authToken,
		})
	})
	router.POST("/upload", controller2.UploadFile)
	router.POST("/download/:id", controller2.DownloadFile)
	router.Run(":8080")
}
