package lib

import (
	"errors"

	"github.com/lokalise/go-lokalise-api/v3"
)

func FetchKeyContent(lokaliseProjectId string, keyName string) (map[string]string, error) {
	if len(lokaliseProjectId) == 0 {
		return nil, errors.New("empty lokaliseProjectId")
	}

	if len(keyName) == 0 {
		return nil, errors.New("empty keyName")
	}

	apiKey, err := getLokaliseApiKey()

	if err != nil {
		return nil, err
	}

	client, err := lokalise.New(apiKey)
	if err != nil {
		return nil, err
	}

	response, err := client.Keys().WithListOptions(lokalise.KeyListOptions{IncludeTranslations: 1, FilterKeys: keyName}).List(lokaliseProjectId)
	if err != nil {
		return nil, err
	}

	if len(response.Keys) == 0 {
		return nil, errors.New("key not found")
	}

	if len(response.Keys) > 1 {
		return nil, errors.New("more than one keys returned")
	}

	var result = map[string]string{}
	for _, translation := range response.Keys[0].Translations {
		result[translation.LanguageISO] = translation.Translation
	}

	if len(result) == 0 {
		return nil, errors.New("no translations found")
	}

	return result, nil
}

// private
func getLokaliseApiKey() (string, error) {
	return FetchEnvVar("LOKALISE_API_KEY")
}
