package services

import (
	"github.com/ddefrancesco/go-microservices/src/api/client/restclient"
	"github.com/ddefrancesco/go-microservices/src/api/domain/repositories"
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

func TestCreateRepoInvalidName(t *testing.T) {
	request := repositories.CreateRepoRequest{}

	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t,"nome repository non valido", err.Message())

}
func TestCreateRepoGithubError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   &http.Response{
			StatusCode:       http.StatusUnauthorized,
			Body:             ioutil.NopCloser( strings.NewReader(`{"message": "Requires authentication", "documentation_url": "https://developer.github.com/v3/repos/#create"}`)),
		},

	})
	//Se passo la request vuota sarà un HTTP400 BadRequest
	//Io voglio HTTP401 Unauthorized quindi...
	request := repositories.CreateRepoRequest{
		Name:        "Test Repo",
	}
	result, err := RepositoryService.CreateRepo(request)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t,"Requires authentication", err.Message())


}
func TestCreateRepoNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response:   &http.Response{
			StatusCode:       http.StatusCreated,
			Body:             ioutil.NopCloser( strings.NewReader(`{"id": 123,"name": "test repo", "owner": { "login": "ddefrancesco"}}`)),
		},

	})
	request := repositories.CreateRepoRequest{
		Name:        "Test Repo",
	}
	result, err := RepositoryService.CreateRepo(request)
	assert.Nil(t, err)
	assert.NotNil(t, result )
	assert.EqualValues(t, 123, result.Id)
	assert.EqualValues(t, "test repo", result.Name)
	assert.EqualValues(t, "ddefrancesco", result.Owner)
}