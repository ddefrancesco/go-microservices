package app

import (
	"github.com/ddefrancesco/go-microservices/src/api/controllers/health"
	"github.com/ddefrancesco/go-microservices/src/api/controllers/repositories"
)

func mapUrls() {
	router.GET("/healthCheck", health.Alive)
	router.POST("/repository", repositories.CreateRepo)
	router.POST("/repositories", repositories.CreateRepos)
}
