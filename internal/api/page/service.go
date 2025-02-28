package page

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rtretter/brain-engine-go/internal/api/auth"
	"github.com/rtretter/brain-engine-go/internal/api/page/dto"
	"github.com/rtretter/brain-engine-go/internal/api/page/model"
	"github.com/rtretter/brain-engine-go/internal/util"
)

type pageService struct {
	authService auth.AuthService
}

func NewPageService(authService auth.AuthService) pageService {
	return pageService{
		authService: authService,
	}
}

func (p *pageService) GetPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundCredentials, err := p.authService.GetUserFromAuthorization(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		pageId := r.PathValue("PAGE_ID")
		owner := r.URL.Query().Get("owner")
		if owner == "" {
			owner = foundCredentials.Username
		}
		page, err := util.LoadPage(pageId, owner)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		responsePage := dto.GetPageResponse{
			ID:         page.ID,
			Title:      page.Title,
			Content:    page.Content,
			OwnerName:  page.OwnerName,
			ModifiedAt: page.ModifiedAt,
			CreatedAt:  page.CreatedAt,
			IsDeleted:  page.IsDeleted,
		}
		json.NewEncoder(w).Encode(responsePage)
	}
}

func (p *pageService) CreatePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundCredentials, err := p.authService.GetUserFromAuthorization(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		var createRequest dto.CreatePage
		err = json.NewDecoder(r.Body).Decode(&createRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newPageId := util.RandomStringAlphaNumerical(32, 32)
		newPage := model.Page{
			ID:         newPageId,
			OwnerName:  foundCredentials.Username,
			CreatedAt:  time.Now(),
			ModifiedAt: time.Now(),
			Title:      createRequest.Title,
			Content:    createRequest.Content,
			IsDeleted:  false,
		}
		if err := util.SavePage(newPage); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newPage)
	}
}
