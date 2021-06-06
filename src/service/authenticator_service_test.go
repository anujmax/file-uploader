package service

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("AUTH_TOKEN", "test_token")
}

func TestAuthenticateShouldGetToken(t *testing.T) {
	//orig := os.Getenv("AUTH_TOKEN")
	//os.Setenv("AUTH_TOKEN", "test_token")
	token, err := GetAuthToken()

	if token != "test_token" {
		t.Errorf("token = %s; want test_token", token)
	}

	if err != nil {
		t.Errorf("error = %s; want nil", err)
	}
	//t.Cleanup(func() { os.Setenv("AUTH_TOKEN", orig) })
}

/*func (a *AuthenticatorServiceTest) TestAuthenticateShouldGetError(t *testing.T) {
	orig := os.Getenv("AUTH_TOKEN")
	os.Setenv("AUTH_TOKEN", "")
	token, err := GetAuthToken()
	if err == nil {
		t.Errorf("error = %s; want not nil", "nil")
	}
	if token != "" {
		t.Errorf("token = %s; want empty ", token)
	}
	t.Cleanup(func() { os.Setenv("AUTH_TOKEN", orig) })
}*/
