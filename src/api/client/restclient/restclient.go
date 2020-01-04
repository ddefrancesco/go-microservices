package restclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	mocksEnabled = false
	mocks        = make(map[string]*Mock)
)

type Mock struct {
	Url        string
	HttpMethod string
	Response   *http.Response
	Err        error
}

func StartMockups() {
	mocksEnabled = true
}
func StopMockups() {
	mocksEnabled = false
}

func GetMockId(httpMethod string, url string) string {
	return fmt.Sprintf("%s_%s", httpMethod,url)
}

func AddMockup(mock *Mock) map[string]*Mock {

	mocks[GetMockId(mock.HttpMethod,mock.Url)] = mock
	return mocks
}

func FlushMockups()  {
	mocks = make(map[string]*Mock)
}

func Post(url string, body interface{}, headers http.Header) (*http.Response, error) {

	if mocksEnabled {
		mock := mocks[GetMockId(http.MethodPost,url)]
		if mock == nil {
			return nil, errors.New("Mockup per la url non trovato")
		}
		return mock.Response, mock.Err

	}

	jsonBytes, err := json.Marshal(body)
	log.Printf("restclient.go: %s",jsonBytes)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	request.Header = headers

	client := http.Client{}

	return client.Do(request)

}
