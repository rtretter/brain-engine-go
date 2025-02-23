package api

import (
	"log"
	"net/http"

	"github.com/rtretter/brain-engine-go/internal/api/auth"
	"github.com/rtretter/brain-engine-go/internal/util"
)

func SetupRoutes() {
	credentials, err := util.LoadCredentials()

	if err != nil {
		log.Fatalf("Failed loading credentials: %+v", err)
	}

	authService := auth.NewAuthService(credentials)

	httpMux := http.NewServeMux()

	httpMux.HandleFunc("GET /auth", AuthMiddleware(authService.GetAuth(), authService))

	http.ListenAndServe(":8080", httpMux)
}
