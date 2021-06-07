package main

import (
	"github.com/anujmax/file-uploader/src/controller"
	"github.com/anujmax/file-uploader/src/repository"
	"github.com/anujmax/file-uploader/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	authToken := os.Getenv("AUTH_TOKEN")
	service.AuthenticationService.Initialize(authToken)
	username := os.Getenv("mysql_username")
	password := os.Getenv("mysql_password")
	host := os.Getenv("mysql_host")
	schema := os.Getenv("mysql_schema")
	repository.FileMetaRepo.Initialize(username, password, host, schema)
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"auth_token": authToken,
		})
	})
	router.POST("/upload", controller.UploadFile)
	router.GET("/download/:id", controller.DownloadFile)
	router.Run(":8080")
}
