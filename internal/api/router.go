package api

import (
	"log"
	"net/http"

	"github.com/rtretter/brain-engine-go/internal/api/auth"
	"github.com/rtretter/brain-engine-go/internal/api/page"
	"github.com/rtretter/brain-engine-go/internal/util"
)

func SetupRoutes() {
	credentials, err := util.LoadCredentials()

	if err != nil {
		log.Fatalf("Failed loading credentials: %+v", err)
	}

	authService := auth.NewAuthService(credentials)
	pageService := page.NewPageService(authService)

	httpMux := http.NewServeMux()

	httpMux.HandleFunc("GET /auth", AuthMiddleware(authService.GetAuth(), authService))

	httpMux.HandleFunc("GET /pages", AuthMiddleware(pageService.QueryPages(), authService))
	httpMux.HandleFunc("POST /pages", AuthMiddleware(pageService.CreatePage(), authService))
	httpMux.HandleFunc("GET /page/{PAGE_ID}", ValidPageIDMiddleware(AuthMiddleware(pageService.GetPage(), authService)))

	http.ListenAndServe(":8080", httpMux)
}
