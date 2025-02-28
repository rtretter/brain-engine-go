package api

import (
	"net/http"
	"regexp"

	"github.com/rtretter/brain-engine-go/internal/api/auth"
)

func AuthMiddleware(next http.Handler, authService auth.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := authService.GetUserFromAuthorization(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func ValidPageIDMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathParam := r.PathValue("PAGE_ID")
		validPageId := isValidPageID(pathParam)
		if !validPageId {
			http.Error(w, "Invalid value for page ID", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func isValidPageID(pageId string) bool {
	isValid, err := regexp.MatchString("^([A-Z]|[a-z]|[0-9])+$", pageId)
	if err != nil {
		panic("Regex failed to check if page ID is valid!")
	}
	return isValid
}
