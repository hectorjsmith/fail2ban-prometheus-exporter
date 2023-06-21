package auth

import (
	"fmt"
	"net/http"
)

func NewBasicAuthProvider(username, password string) AuthProvider {
	return &basicAuthProvider{
		hashedAuth: encodeBasicAuth(username, password),
	}
}

type basicAuthProvider struct {
	hashedAuth string
}

func (p *basicAuthProvider) IsAllowed(request *http.Request) bool {
	username, password, ok := request.BasicAuth()
	if !ok {
		return false
	}
	requestAuth := encodeBasicAuth(username, password)
	return p.hashedAuth == requestAuth
}

func encodeBasicAuth(username, password string) string {
	return HashString(fmt.Sprintf("%s:%s", username, password))
}
