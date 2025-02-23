package api

import (
	"net/http"

	"github.com/rtretter/brain-engine-go/internal/api/auth"
)

func AuthMiddleware(next http.Handler, authService auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := authService.GetUserFromAuthorization(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}
