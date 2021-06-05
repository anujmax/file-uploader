package main

import (
	"github.com/anujmax/file-uploader/controller"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"auth_token": "askjlfjkasd978987897asdf",
		})
	})
	router.POST("/upload", controller.UploadFile)
	router.POST("/download/:id", controller.DownloadFile)
	router.Run(":8080")
}
