package util

import (
	"encoding/json"
	"log"
	"os"

	"github.com/rtretter/brain-engine-go/internal/api/auth/model"
)

const dataPath = "brain"
const credentialFile = dataPath + "/credentials.json"

func LoadCredentials() (credentials *[]model.Credentials, err error) {
	loadedCredentials, err := loadJson[[]model.Credentials](credentialFile)
	if err != nil {
		loadedCredentials, err = generateCredentials()
	}
	return loadedCredentials, err
}

func generateCredentials() (*[]model.Credentials, error) {
	generatedCredentials := []model.Credentials{
		{
			Username: "admin",
			Token:    randomStringDefaultCharset(32, 42),
		},
	}
	err := saveJson(credentialFile, generatedCredentials)
	return &generatedCredentials, err
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
