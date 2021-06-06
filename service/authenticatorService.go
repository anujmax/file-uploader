package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"os"
)

var (
	authToken = os.Getenv("AUTH_TOKEN")
)

func Authenticate(c *gin.Context) error {
	authToken := c.Request.FormValue("token")
	if authToken == authToken {
		return nil
	} else {
		return errors.New("invalid auth token")
	}
}

func GetAuthToken() (string, error) {
	if len(authToken) > 0 {
		return authToken, nil
	} else {
		return "", errors.New("authentication token not present in env variable")
	}
}
