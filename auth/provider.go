package auth

import (
	"net/http"
)

type AuthProvider interface {
	IsAllowed(*http.Request) bool
}

func NewEmptyAuthProvider() AuthProvider {
	return &emptyAuthProvider{}
}

type emptyAuthProvider struct {
}

func (p *emptyAuthProvider) IsAllowed(request *http.Request) bool {
	return true
}

type compositeAuthProvider struct {
	providers []AuthProvider
}

func (p *compositeAuthProvider) IsAllowed(request *http.Request) bool {
	for i := 0; i < len(p.providers); i++ {
		provider := p.providers[i]
		if provider.IsAllowed(request) {
			return true
		}
	}
	return false
}
