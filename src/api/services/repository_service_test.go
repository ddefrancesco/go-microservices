package services

import (
	"github.com/ddefrancesco/go-microservices/src/api/client/restclient"
	"github.com/ddefrancesco/go-microservices/src/api/domain/repositories"
	"github.com/ddefrancesco/go-microservices/src/api/utils/errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestMain(m *testing.M) {
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidName(t *testing.T) {
	request := repositories.CreateRepoRequest{}

	result, err := RepositoryService.CreateRepo(request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusBadRequest, err.Status())
	assert.EqualValues(t, "nome repository non valido", err.Message())

}
func TestCreateRepoGithubError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication", "documentation_url": "https://developer.github.com/v3/repos/#create"}`)),
		},
	})
	//Se passo la request vuota sarà un HTTP400 BadRequest
	//Io voglio HTTP401 Unauthorized quindi...
	request := repositories.CreateRepoRequest{
		Name: "Test Repo",
	}
	result, err := RepositoryService.CreateRepo(request)
	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusUnauthorized, err.Status())
	assert.EqualValues(t, "Requires authentication", err.Message())

}
func TestCreateRepoNoError(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123,"name": "test repo", "owner": { "login": "ddefrancesco"}}`)),
		},
	})
	request := repositories.CreateRepoRequest{
		Name: "Test Repo",
	}
	result, err := RepositoryService.CreateRepo(request)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, 123, result.Id)
	assert.EqualValues(t, "test repo", result.Name)
	assert.EqualValues(t, "ddefrancesco", result.Owner)
}

func TestReposService_CreateRepos_InvalidRequest(t *testing.T) {
	request := repositories.CreateRepoRequest{}

	output := make(chan repositories.CreateRepositoriesResult)

	service := reposService{}

	go service.CreateRepoConcurrent(request, output)

	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
	assert.EqualValues(t, "nome repository non valido", result.Error.Message())

}

func TestReposService_CreateRepos_ErrorFromGitlab(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusUnauthorized,
			Body:       ioutil.NopCloser(strings.NewReader(`{"message": "Requires authentication", "documentation_url": "https://developer.github.com/v3/repos/#create"}`)),
		},
	})

	request := repositories.CreateRepoRequest{Name: "testing"}
	output := make(chan repositories.CreateRepositoriesResult)

	service := reposService{}

	go service.CreateRepoConcurrent(request, output)

	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Response)
	assert.NotNil(t, result.Error)
	assert.EqualValues(t, http.StatusUnauthorized, result.Error.Status())
	assert.EqualValues(t, "Requires authentication", result.Error.Message())

}

func TestReposService_CreateRepos_NoError(t *testing.T) {

	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123,"name": "test repo", "owner": { "login": "ddefrancesco"}}`)),
		},
	})

	request := repositories.CreateRepoRequest{Name: "testing"}
	output := make(chan repositories.CreateRepositoriesResult)

	service := reposService{}

	go service.CreateRepoConcurrent(request, output)

	result := <-output
	assert.NotNil(t, result)
	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Response)
	assert.EqualValues(t, 123, result.Response.Id)
	assert.EqualValues(t, "test repo", result.Response.Name)
	assert.EqualValues(t, "ddefrancesco", result.Response.Owner)

}
func TestReposService_CreateRepo_HandleRepoResults(t *testing.T) {
	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)

	var wg sync.WaitGroup

	service := reposService{}
	go service.handleRepoResults(&wg, input, output)

	wg.Add(1)
	go func() {
		input <- repositories.CreateRepositoriesResult{

			Error: errors.NewBadRequestError("invalid repo request"),
		}
	}()
	wg.Wait()
	close(input)
	result := <-output
	assert.NotNil(t, result)
	assert.EqualValues(t, 1, len(result.Results))
	assert.NotNil(t, result.Results[0].Error)
	assert.Nil(t, result.Results[0].Response)
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repo request", result.Results[0].Error.Message())
}

func TestReposService_CreateRepos_Invalid_Requests(t *testing.T) {
	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "   "},
	}

	result, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusBadRequest, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	assert.Nil(t, result.Results[0].Response)
	assert.EqualValues(t, http.StatusBadRequest, result.Results[0].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[0].Error.Message())

	assert.Nil(t, result.Results[1].Response)
	assert.EqualValues(t, http.StatusBadRequest, result.Results[1].Error.Status())
	assert.EqualValues(t, "invalid repository name", result.Results[1].Error.Message())

}

func TestCreateReposOneSuccessOneFail(t *testing.T) {
	restclient.FlushMockups()
	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode: http.StatusCreated,
			Body:       ioutil.NopCloser(strings.NewReader(`{"id": 123, "name": "testing", "owner": {"login": "federicoleon"}}`)),
		},
	})

	requests := []repositories.CreateRepoRequest{
		{},
		{Name: "testing"},
	}

	result, err := RepositoryService.CreateRepos(requests)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.EqualValues(t, http.StatusPartialContent, result.StatusCode)
	assert.EqualValues(t, 2, len(result.Results))

	for _, result := range result.Results {
		if result.Error != nil {
			assert.EqualValues(t, http.StatusBadRequest, result.Error.Status())
			assert.EqualValues(t, "invalid repository name", result.Error.Message())
			continue
		}

		assert.EqualValues(t, 123, result.Response.Id)
		assert.EqualValues(t, "testing", result.Response.Name)
		assert.EqualValues(t, "federicoleon", result.Response.Owner)
	}
}
