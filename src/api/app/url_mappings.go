package app

import (
	"github.com/ddefrancesco/go-microservices/src/api/controllers/health"
	"github.com/ddefrancesco/go-microservices/src/api/controllers/repositories"
)

func mapUrls()  {
	router.GET("/healthCheck",health.Alive)
	router.POST("/repositories", repositories.CreateRepo)
}
