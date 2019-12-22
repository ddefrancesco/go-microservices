package services

import (
	"github.com/ddefrancesco/go-microservices/src/api/config"
	"github.com/ddefrancesco/go-microservices/src/api/domain/github"
	"github.com/ddefrancesco/go-microservices/src/api/domain/repositories"
	"github.com/ddefrancesco/go-microservices/src/api/providers/github_provider"
	"github.com/ddefrancesco/go-microservices/src/api/utils/errors"
	"strings"
)

type reposService struct {}

type reposServiceInterface interface {
	CreateRepo(inputRequest repositories.CreateRepoRequest) (*repositories.CreateRepoResponse,errors.ApiError)
}

var (
	RepositoryService reposServiceInterface
)

func init()  {

	RepositoryService = &reposService{}
}

func (m *reposService) CreateRepo(inputRequest repositories.CreateRepoRequest) (*repositories.CreateRepoResponse, errors.ApiError){
	inputRequest.Name = strings.TrimSpace(inputRequest.Name)
	if inputRequest.Name == ""{
		return nil,errors.NewBadRequestError("nome repository non valido")
	}

	request := github.CreateRepoRequest{
		Name:        inputRequest.Name,
		Description: inputRequest.Description,
		Private:     false,

	}

	response, err := github_provider.CreateRepo(config.GetGithubAccessToken(),request)
	if err != nil {

		return nil, errors.NewApiError(err.StatusCode,err.Message)
	}

	result := repositories.CreateRepoResponse{
		Id:    response.Id,
		Owner: response.Owner.Login,
		Name:  response.Name,
	}

	return &result,nil
}