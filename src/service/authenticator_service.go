package service

import (
	"errors"
	"os"
)

var (
	authenticationToken = os.Getenv("AUTH_TOKEN")
)

func Authenticate(authToken string) error {
	if authToken == authenticationToken {
		return nil
	} else {
		return errors.New("invalid auth token")
	}
}

func GetAuthToken() (string, error) {
	if len(authenticationToken) > 0 {
		return authenticationToken, nil
	} else {
		return "", errors.New("authentication token not present in env variable")
	}
}
