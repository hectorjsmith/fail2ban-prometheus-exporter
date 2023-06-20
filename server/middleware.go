package server

import (
	"net/http"
)

type BasicAuthProvider interface {
	Enabled() bool
	DoesBasicAuthMatch(username, password string) bool
}

func BasicAuthMiddleware(handlerFunc http.HandlerFunc, basicAuthProvider BasicAuthProvider) http.HandlerFunc {
	if basicAuthProvider.Enabled() {
		return func(w http.ResponseWriter, r *http.Request) {
			if doesBasicAuthMatch(r, basicAuthProvider) {
				handlerFunc.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		}
	}
	return handlerFunc
}

func doesBasicAuthMatch(r *http.Request, basicAuthProvider BasicAuthProvider) bool {
	rawUsername, rawPassword, ok := r.BasicAuth()
	if ok {
		return basicAuthProvider.DoesBasicAuthMatch(rawUsername, rawPassword)
	}
	return false
}
