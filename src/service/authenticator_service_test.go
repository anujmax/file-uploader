package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticateSuccess(t *testing.T) {
	AuthenticationService.Initialize("test_token")
	err := AuthenticationService.Authenticate("test_token")
	assert.Nil(t, err)
}

func TestAuthenticateFailure(t *testing.T) {
	AuthenticationService.Initialize("test_token")
	err := AuthenticationService.Authenticate("fake_token")
	assert.Equal(t, err, errors.New("invalid auth token"))
}
