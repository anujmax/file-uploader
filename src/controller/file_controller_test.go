package controller

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestUploadFile(t *testing.T) {

	var c *gin.Context
	UploadFile(c)
	if c == nil {
		t.Error("Expected Register User to throw and error got nil")
	}
}