package auth

import "net/http"

func NewEmptyAuthProvider() AuthProvider {
	return &emptyAuthProvider{}
}

type emptyAuthProvider struct {
}

func (p *emptyAuthProvider) IsAllowed(request *http.Request) bool {
	return true
}
