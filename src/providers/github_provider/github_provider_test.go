package github_provider

import (
	"errors"
	"github.com/ddefrancesco/go-microservices/src/client/restclient"
	"github.com/ddefrancesco/go-microservices/src/domain/github"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M)  {

	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestGetAuthorizationHeader(t *testing.T)  {
	access_token := "123abc"
	auth_header := getAuthorizationHeader(access_token)
	assert.EqualValues(t, "token 123abc" , auth_header)
}

func TestCreateRepoErrorRestclient(t *testing.T){
	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   nil,
		Err:        errors.New("invalid restclient repo request"),
	})
	response, err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t, http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t, "invalid restclient repo request",err.Message)
}

func TestCreateRepoInvalidResponseBody(t *testing.T){
	restclient.FlushMockups()
	invalidBody, _ := os.Open("abc")
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode:       http.StatusCreated,
			Body:             invalidBody,
		},
		Err: nil,
	})
	response, err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t, http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t, "invalid response body",err.Message)
}

func TestCreateRepoInvalidErrorInterface(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode:       http.StatusUnauthorized,
			Body:             ioutil.NopCloser( strings.NewReader(`{"message": 1}`)),
		},
		Err: nil,
	})
	response, err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t, http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t, "invalid json body",err.Message)
}

func TestCreateRepoUnauthorized(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode:       http.StatusUnauthorized,
			Body:             ioutil.NopCloser( strings.NewReader(`{"message": "Requires authentication", "documentation_url": "https://developer.github.com/v3/repos/#create"}`)),
		},
		Err: nil,
	})
	response, err := CreateRepo("",github.CreateRepoRequest{})
	assert.Nil(t,response)
	assert.NotNil(t,err)
	assert.EqualValues(t, http.StatusUnauthorized,err.StatusCode)
	assert.EqualValues(t, "Requires authentication",err.Message)
}

func TestCreateRepoInvalidSuccessResponse(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode:       http.StatusCreated,
			Body:             ioutil.NopCloser( strings.NewReader(`{"id": "123"}`)),
		},

	})
	response, err := CreateRepo("",github.CreateRepoRequest{})
	assert.NotNil(t,err)
	assert.Nil(t,response)
	assert.EqualValues(t, http.StatusInternalServerError,err.StatusCode)
	assert.EqualValues(t, "error unmarshalling create repo response",err.Message)
}

func TestCreateRepoNoErrorResponse(t *testing.T){
	restclient.FlushMockups()

	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode:       http.StatusCreated,
			Body:             ioutil.NopCloser( strings.NewReader(`{"id": 1234567, "name": "Hello-World","full_name": "ddefrancesco/Hello-World"}`)),
		},

	})
	response, err := CreateRepo("",github.CreateRepoRequest{})
	assert.NotNil(t,response)
	assert.Nil(t,err)

	assert.EqualValues(t, 1234567,response.Id)
	assert.EqualValues(t, "Hello-World",response.Name)
	assert.EqualValues(t, "ddefrancesco/Hello-World",response.FullName)
}