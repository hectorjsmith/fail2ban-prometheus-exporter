package auth

import (
	"net/http"
)

type AuthProvider interface {
	IsAllowed(*http.Request) bool
}
