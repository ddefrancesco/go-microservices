package services

import (
	"github.com/ddefrancesco/go-microservices/src/api/config"
	"github.com/ddefrancesco/go-microservices/src/api/domain/github"
	"github.com/ddefrancesco/go-microservices/src/api/domain/repositories"
	"github.com/ddefrancesco/go-microservices/src/api/providers/github_provider"
	"github.com/ddefrancesco/go-microservices/src/api/utils/errors"
	"net/http"
	"sync"
)

type reposService struct{}

type reposServiceInterface interface {
	CreateRepo(inputRequest repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError)
	CreateRepos(inputRequest []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init() {

	RepositoryService = &reposService{}
}

func (m *reposService) CreateRepo(inputRequest repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError) {

	if err := inputRequest.Validate(); err != nil {
		return nil, err
	}

	request := github.CreateRepoRequest{
		Name:        inputRequest.Name,
		Description: inputRequest.Description,
		Private:     false,
	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(), request)
	if err != nil {

		return nil, errors.NewApiError(err.StatusCode, err.Message)
	}

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Name:  response.Name,
		Owner: response.Owner.Login,
	}

	return &result, nil
}

func (m *reposService) CreateRepos(inputRequest []repositories.CreateRepoRequest) (repositories.CreateReposResponse, errors.ApiError) {

	input := make(chan repositories.CreateRepositoriesResult)
	output := make(chan repositories.CreateReposResponse)

	defer close(output)

	var wg sync.WaitGroup //Inits wait group

	go m.handleRepoResults(&wg, input, output)

	for _, current := range inputRequest {
		wg.Add(1) //adds 1 @implicit counter

		go m.CreateRepoConcurrent(current, input)
	}
	wg.Wait()    //Suspends execution until implicit counter doesn't reach zero
	close(input) // close the input channel
	result := <-output

	successOutcomes := 0

	for _, current := range result.Results {
		if current.Response != nil {
			successOutcomes++
		}
	}

	if successOutcomes == 0 {
		result.StatusCode = result.Results[0].Error.Status() // if we failed all the request then take first CreateRepoResponse
		// -> Error -> Status (github error status code)

	} else if successOutcomes == len(inputRequest) {
		result.StatusCode = http.StatusCreated
	} else {
		result.StatusCode = http.StatusPartialContent
	}
	return result, nil
}

func (m *reposService) handleRepoResults(wg *sync.WaitGroup, input chan repositories.CreateRepositoriesResult, output chan repositories.CreateReposResponse) {
	var results repositories.CreateReposResponse

	for incomingEvent := range input {
		repositoriesResult := repositories.CreateRepositoriesResult{
			Response: incomingEvent.Response,
			Error:    incomingEvent.Error,
		}
		results.Results = append(results.Results, repositoriesResult)
		wg.Done() //Subtracts 1 @implicit counter
	}

	output <- results
}

func (m *reposService) CreateRepoConcurrent(input repositories.CreateRepoRequest, output chan repositories.CreateRepositoriesResult) {

	if err := input.Validate(); err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
	}
	result, err := m.CreateRepo(input)
	if err != nil {
		output <- repositories.CreateRepositoriesResult{Error: err}
		return
	}
	output <- repositories.CreateRepositoriesResult{Response: result}

}
