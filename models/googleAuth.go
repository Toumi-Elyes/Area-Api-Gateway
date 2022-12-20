package models

import (
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
)

type AuthModel struct {
	Code string `json:"code"`
}

func GetCredentialConfig(filename string) (*oauth2.Config, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(fileContent, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
		return nil, err
	}
	return config, nil
}
