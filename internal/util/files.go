package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	authModel "github.com/rtretter/brain-engine-go/internal/api/auth/model"
	pageModel "github.com/rtretter/brain-engine-go/internal/api/page/model"
)

const dataPath = "brain/"
const pagesPath = dataPath + "pages/"    // DATA_PATH/pages/
const ownedPagesPath = pagesPath + "%s/" // DATA_PATH/pages/USERNAME/
const pagePath = ownedPagesPath + "%s/"  // DATA_PATH/pages/USERNAME/PAGE_NAME/
const pageFile = pagePath + "page.json"
const credentialFile = dataPath + "credentials.json"

func LoadCredentials() (credentials *[]authModel.Credentials, err error) {
	loadedCredentials, err := loadJson[[]authModel.Credentials](credentialFile)
	if err != nil {
		loadedCredentials, err = generateCredentials()
	}
	return loadedCredentials, err
}

func generateCredentials() (*[]authModel.Credentials, error) {
	generatedCredentials := []authModel.Credentials{
		{
			Username: "admin",
			Token:    RandomStringDefaultCharset(32, 42),
		},
	}
	err := saveJson(credentialFile, generatedCredentials)
	return &generatedCredentials, err
}

func QueryAllPages(query string, includeDeleted bool) (*[]pageModel.Page, error) {
	ownerDirs, err := os.ReadDir(pagesPath)
	if err != nil {
		return nil, err
	}
	var response []pageModel.Page = make([]pageModel.Page, 0)
	for _, owner := range ownerDirs {
		ownedPages, err := QueryOwnPages(query, owner.Name(), includeDeleted)
		if err != nil {
			continue
		}
		response = append(response, *ownedPages...)
	}
	return &response, nil
}

func QueryOwnPages(query, owner string, includeDeleted bool) (*[]pageModel.Page, error) {
	actualPagesDir := fmt.Sprintf(ownedPagesPath, owner)
	pageDirs, err := os.ReadDir(actualPagesDir)
	if err != nil {
		return nil, err
	}
	var response []pageModel.Page = make([]pageModel.Page, 0)
	appendMatchingPages(pageDirs, query, owner, includeDeleted, &response)
	return &response, nil
}

func appendMatchingPages(pageDirs []os.DirEntry, query, owner string, includeDeleted bool, pageResponse *[]pageModel.Page) {
	for _, e := range pageDirs {
		if !e.IsDir() {
			continue
		}
		page, err := LoadPage(e.Name(), owner)
		if err != nil {
			continue
		}
		if page.IsDeleted && !includeDeleted {
			continue
		}
		if strings.Contains(strings.ToLower(page.Title), strings.ToLower(query)) {
			*pageResponse = append(*pageResponse, *page)
			continue
		}
		if strings.Contains(strings.ToLower(page.Content), strings.ToLower(query)) {
			*pageResponse = append(*pageResponse, *page)
			continue
		}
	}

}

func LoadPage(pageId, owner string) (*pageModel.Page, error) {
	actualPageFile := fmt.Sprintf(pageFile, owner, pageId)
	if _, err := os.Stat(actualPageFile); errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	page, err := loadJson[pageModel.Page](actualPageFile)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func SavePage(page pageModel.Page) error {
	actualPageDir := fmt.Sprintf(pagePath, page.OwnerName, page.ID)
	actualPageFile := fmt.Sprintf(pageFile, page.OwnerName, page.ID)
	if _, err := os.Stat(actualPageDir); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(actualPageDir, 0766); err != nil {
			return err
		}
	}
	err := saveJson(actualPageFile, page)
	if err != nil {
		return err
	}
	return nil
}

func saveJson(file string, data any) error {
	jsonData, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if err := os.WriteFile(file, jsonData, 0666); err != nil {
		return err
	}

	return nil
}

func loadJson[dataType any](file string) (data *dataType, err error) {
	content, err := os.ReadFile(file)

	if err != nil {
		if err := os.MkdirAll(dataPath, 0766); err != nil {
			log.Fatalf("Unable to create data dir(%s): %s", dataPath, err)
		}
		return nil, err
	}

	var response dataType
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
