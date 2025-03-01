package page

import (
	"encoding/json"
	"net/http"
	"strconv"
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

func (p *pageService) QueryPages() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundCredentials, _ := p.authService.GetUserFromAuthorization(r.Header.Get("Authorization"))
		query := r.URL.Query().Get("query")
		includeDeleted, err := strconv.ParseBool(r.URL.Query().Get("includeDeleted"))
		if err != nil {
			includeDeleted = false
		}
		includeUnowned, err := strconv.ParseBool(r.URL.Query().Get("includeUnowned"))
		if err != nil {
			includeUnowned = false
		}
		err = nil
		var pages *[]model.Page
		if includeUnowned {
			pages, err = util.QueryAllPages(query, includeDeleted)
		} else {
			pages, err = util.QueryOwnPages(query, foundCredentials.Username, includeDeleted)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(pages)
	}
}

func (p *pageService) GetPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundCredentials, _ := p.authService.GetUserFromAuthorization(r.Header.Get("Authorization"))
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
		json.NewEncoder(w).Encode(page)
	}
}

func (p *pageService) CreatePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		foundCredentials, _ := p.authService.GetUserFromAuthorization(r.Header.Get("Authorization"))
		var createRequest dto.CreatePage
		err := json.NewDecoder(r.Body).Decode(&createRequest)
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
