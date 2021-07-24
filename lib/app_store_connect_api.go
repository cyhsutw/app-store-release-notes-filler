package lib

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AppVersion struct {
	Id    string
	State string
}

type AppVersionLocalization struct {
	Id           string
	Locale       string
	ReleaseNotes string
}

func UpdateVersionLocalization(model AppVersionLocalization) (AppVersionLocalization, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("https://api.appstoreconnect.apple.com/v1/appStoreVersionLocalizations/%s", model.Id))
	if err != nil {
		message := fmt.Sprintf("error building request url: %v", err)
		return model, errors.New(message)
	}

	requestPayload := gin.H{
		"data": gin.H{
			"id": model.Id,
			"attributes": gin.H{
				"whatsNew": model.ReleaseNotes,
			},
			"type": "appStoreVersionLocalizations",
		},
	}

	requestPayloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		message := fmt.Sprintf("error serialize payload to json bytes: %v", err)
		return model, errors.New(message)
	}

	request, err := http.NewRequest("PATCH", requestUrl.String(), bytes.NewBuffer(requestPayloadBytes))
	if err != nil {
		message := fmt.Sprintf("error creating request: %v", err)
		return model, errors.New(message)
	}

	request.Header.Set("Content-Type", "application/json")

	statusCode, payload, err := sendRequest(request)
	if err != nil {
		message := fmt.Sprintf("error sending request: %v", err)
		return model, errors.New(message)
	}

	if statusCode != http.StatusOK {
		message := fmt.Sprintf("unexpected response code: %d", statusCode)
		return model, errors.New(message)
	}

	newModel, err := createAppVersionLocalization(payload.Get("data").MustMap())

	if err != nil {
		message := fmt.Sprintf("error creating AppVersionLocalization: %v", err)
		return model, errors.New(message)
	}

	return newModel, nil
}

func FetchVersionLocalizations(versionId string) (map[string]AppVersionLocalization, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("https://api.appstoreconnect.apple.com/v1/appStoreVersions/%s/appStoreVersionLocalizations", versionId))
	if err != nil {
		message := fmt.Sprintf("error building request url: %v", err)
		return nil, errors.New(message)
	}

	queries := requestUrl.Query()

	includedFields := "whatsNew,locale"
	queries.Set("fields[appStoreVersionLocalizations]", includedFields)

	requestUrl.RawQuery = queries.Encode()

	request, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		message := fmt.Sprintf("error creating request: %v", err)
		return nil, errors.New(message)
	}

	statusCode, payload, err := sendRequest(request)

	if err != nil {
		message := fmt.Sprintf("error sending request: %v", err)
		return nil, errors.New(message)
	}

	if statusCode != http.StatusOK {
		message := fmt.Sprintf("unexpected response code: %d", statusCode)
		return nil, errors.New(message)
	}

	entries, err := payload.Get("data").Array()
	if err != nil {
		return nil, errors.New("no localizations in `data`")
	}

	var results = map[string]AppVersionLocalization{}
	for _, entry := range entries {
		localization, err := createAppVersionLocalization(entry)
		if err != nil {
			message := fmt.Sprintf("error creating AppVersionLocalization: %v", err)
			return nil, errors.New(message)
		}
		results[localization.Locale] = localization
	}

	return results, nil
}

func FetchEditableVersion(appId string) (AppVersion, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("https://api.appstoreconnect.apple.com/v1/apps/%s/appStoreVersions", appId))
	if err != nil {
		message := fmt.Sprintf("error building request url: %v", err)
		return AppVersion{}, errors.New(message)
	}

	queries := requestUrl.Query()

	states := "DEVELOPER_REMOVED_FROM_SALE,DEVELOPER_REJECTED,INVALID_BINARY,METADATA_REJECTED,PREPARE_FOR_SUBMISSION,REJECTED,REMOVED_FROM_SALE,WAITING_FOR_REVIEW"
	queries.Set("filter[appStoreState]", states)

	includedFields := "appStoreState"
	queries.Set("fields[appStoreVersions]", includedFields)

	queries.Set("limit", "1")

	requestUrl.RawQuery = queries.Encode()

	request, err := http.NewRequest("GET", requestUrl.String(), nil)
	if err != nil {
		message := fmt.Sprintf("error creating request: %v", err)
		return AppVersion{}, errors.New(message)
	}

	statusCode, payload, err := sendRequest(request)

	if err != nil {
		message := fmt.Sprintf("error sending request: %v", err)
		return AppVersion{}, errors.New(message)
	}

	if statusCode != http.StatusOK {
		message := fmt.Sprintf("unexpected response code: %d", statusCode)
		return AppVersion{}, errors.New(message)
	}

	version := payload.Get("data").GetIndex(0)
	if len(version.MustMap()) == 0 {
		return AppVersion{}, errors.New("no version is editable")
	}

	id, err := version.Get("id").String()
	if err != nil {
		message := fmt.Sprintf("cannot get `id` from `data` entry: %v", err)
		return AppVersion{}, errors.New(message)
	}

	state, err := version.GetPath("attributes", "appStoreState").String()
	if err != nil {
		message := fmt.Sprintf("cannot get `attributes.appStoreState` from `data` entry: %v", err)
		return AppVersion{}, errors.New(message)
	}

	return AppVersion{Id: id, State: state}, nil
}

// private
type signingKey interface {
	value() *ecdsa.PrivateKey
}

type sharedSigningKey struct {
	key      *ecdsa.PrivateKey
	mutex    sync.RWMutex
	isLoaded bool
}

func (k *sharedSigningKey) value() *ecdsa.PrivateKey {
	k.mutex.RLock()
	defer k.mutex.RLocker().Unlock()
	return k.key
}

var currentSigningKey = sharedSigningKey{isLoaded: false}

func loadSigningKey() {
	currentSigningKey.mutex.Lock()
	defer currentSigningKey.mutex.Unlock()

	if currentSigningKey.isLoaded {
		return
	}

	currentSigningKey.key = fetchPrivateKey()
	currentSigningKey.isLoaded = true
}

func fetchPrivateKey() *ecdsa.PrivateKey {
	privateKey, err := FetchEnvVar("APP_STORE_CONNECT_API_PRIVATE_KEY")
	if err != nil {
		log.Fatalln("Cannot fetch APP_STORE_CONNECT_API_PRIVATE_KEY")
	}

	block, _ := pem.Decode([]byte(privateKey))
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatalln("Failed to decode PEM block containing private key")
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalln("Failed to create x509 key")
	}

	ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
	if ok == false {
		log.Fatalln("Failed to cast x509 key to ECDSA key")
	}

	return ecdsaPrivateKey
}

func fetchIssuerId() (string, error) {
	return FetchEnvVar("APP_STORE_CONNECT_API_ISSUER_ID")
}

func fetchKeyId() (string, error) {
	return FetchEnvVar("APP_STORE_CONNECT_API_KEY_ID")
}

func generateJwt(issuerId string, keyId string, privateKey signingKey) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodES256,
		jwt.MapClaims{
			"aud": "appstoreconnect-v1",
			"iss": issuerId,
			// valid for 10 minutes
			// (max 20 minutes according to https://developer.apple.com/documentation/appstoreconnectapi/generating_tokens_for_api_requests)
			"exp": time.Now().Add(10 * 60).Unix(),
		})

	token.Header["kid"] = keyId

	signedToken, err := token.SignedString(privateKey.value())
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func sendRequest(request *http.Request) (int, *simplejson.Json, error) {
	client := &http.Client{}

	// generate JWT
	loadSigningKey()

	iss, err := fetchIssuerId()
	if err != nil {
		return -1, nil, err
	}

	kid, err := fetchKeyId()
	if err != nil {
		return -1, nil, err
	}

	token, err := generateJwt(iss, kid, &currentSigningKey)
	if err != nil {
		message := fmt.Sprintf("error generating jwt : %v", err)
		return -1, nil, errors.New(message)
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := client.Do(request)
	if err != nil {
		message := fmt.Sprintf("error getting response: %v", err)
		return response.StatusCode, nil, errors.New(message)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		message := fmt.Sprintf("error reading response body: %v", err)
		return response.StatusCode, nil, errors.New(message)
	}

	if len(bodyBytes) == 0 {
		return response.StatusCode, simplejson.New(), nil
	}

	payload, err := simplejson.NewJson(bodyBytes)
	if err != nil {
		message := fmt.Sprintf("error deserialize response body to json: %v", err)
		return response.StatusCode, nil, errors.New(message)
	}

	return response.StatusCode, payload, nil
}

func createAppVersionLocalization(input interface{}) (AppVersionLocalization, error) {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		message := fmt.Sprintf("error serializing `input` to []byte: %v", err)
		return AppVersionLocalization{}, errors.New(message)
	}

	object, err := simplejson.NewJson(inputBytes)
	if err != nil {
		message := fmt.Sprintf("error deserializing `inputBytes` to simplejson.Json: %v", err)
		return AppVersionLocalization{}, errors.New(message)
	}

	localization := AppVersionLocalization{
		Id:           object.Get("id").MustString(),
		Locale:       object.GetPath("attributes", "locale").MustString(),
		ReleaseNotes: object.GetPath("attributes", "whatsNew").MustString(),
	}

	return localization, nil
}
