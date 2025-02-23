package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/rtretter/brain-engine-go/internal/api/auth/dto"
	"github.com/rtretter/brain-engine-go/internal/api/auth/model"
)

type AuthService struct {
	credentials *[]model.Credentials
}

func NewAuthService(credentials *[]model.Credentials) AuthService {
	return AuthService{credentials: credentials}
}

func (s AuthService) GetUserFromAuthorization(auth string) (*model.Credentials, error) {
	tokens := strings.Split(auth, " ")
	if len(tokens) != 2 {
		return nil, errors.New("Failed splitting bearer token!")
	}
	token := tokens[1]
	var foundCredentials *model.Credentials = nil

	for _, creds := range *s.credentials {
		if creds.Token == token {
			foundCredentials = &creds
		}
	}
	if foundCredentials == nil {
		return nil, errors.New("Unable to find user for token!")
	}
	return foundCredentials, nil
}

func (s AuthService) GetAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundCredentials, err := s.GetUserFromAuthorization(r.Header.Get("Authorization"))
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		response := dto.AuthResponse{
			Username: foundCredentials.Username,
		}
		responseJson, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(responseJson))
	}
}
