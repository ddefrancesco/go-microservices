package github_provider

import (
	"encoding/json"
	"fmt"
	"github.com/ddefrancesco/go-microservices/src/client/restclient"
	"github.com/ddefrancesco/go-microservices/src/domain/github"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	headerAuthorization  = "Authorization"
	headerAuthorizationFmt = "token %s"

	urlCreateRepo = "https://api.github.com/user/repos"
)

func getAuthorizationHeader(access_token string)  string {
	return fmt.Sprintf(headerAuthorizationFmt, access_token)
}

func CreateRepo (access_token string ,request github.CreateRepoRequest) (*github.CreateRepoResponse, *github.GithubErrorResponse){

	headers := http.Header{}
	headers.Set("Authorization",getAuthorizationHeader(access_token))

	response, err := restclient.Post(urlCreateRepo, request,headers)
	if err != nil {
		log.Printf("Error when trying to create github repo: %s",err.Error())
		return nil, &github.GithubErrorResponse{
			StatusCode:       http.StatusInternalServerError,
			Message:             err.Error(),

		}
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, &github.GithubErrorResponse{
			StatusCode:       http.StatusInternalServerError,
			Message:             "invalid response body",

		}
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {
		var errResponse github.GithubErrorResponse
		if err := json.Unmarshal(bytes,&errResponse ); err != nil {
			return nil,  &github.GithubErrorResponse{
				StatusCode:       http.StatusInternalServerError,
				Message:             "invalid json body",

			}
		}
		errResponse.StatusCode = response.StatusCode
		return nil, &errResponse
	}
	var result github.CreateRepoResponse
	if err := json.Unmarshal(bytes,&result ); err != nil {
		return nil,  &github.GithubErrorResponse{
			StatusCode:       http.StatusInternalServerError,
			Message:             "error unmarshalling create repo response",

		}

	}
	return &result, nil
}
