package service

import "errors"

const (
	AUTH_TOKEN = "askjlfjkasd978987897asdf"
)
func Authenticate(authToken string) error {
	if(authToken == AUTH_TOKEN) {
		return nil
	} else {
		return errors.New("Invalid auth token")
	}
}
