package service

import (
	"errors"
)

var (
	AuthenticationService authenticationServiceInterface = &authenticationService{}
)

type authenticationService struct {
	authenticationToken string
}

type authenticationServiceInterface interface {
	Authenticate(string) error
	Initialize(authToken string)
}

func (a *authenticationService) Authenticate(authToken string) error {
	if authToken == a.authenticationToken {
		return nil
	} else {
		return errors.New("invalid auth token")
	}
}

func (a *authenticationService) Initialize(authToken string) {
	a.authenticationToken = authToken
}
