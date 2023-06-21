package server

import (
	"net/http"

	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/auth"
)

func AuthMiddleware(handlerFunc http.HandlerFunc, authProvider auth.AuthProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if authProvider.IsAllowed(r) {
			handlerFunc.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
