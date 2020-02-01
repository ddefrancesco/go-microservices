package repositories

import (
	"github.com/ddefrancesco/go-microservices/src/api/domain/repositories"
	"github.com/ddefrancesco/go-microservices/src/api/services"
	"github.com/ddefrancesco/go-microservices/src/api/utils/errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func CreateRepo(c *gin.Context) {
	log.Printf("repositories_controller.go Request Body: %s", c.Request.Body)
	var request repositories.CreateRepoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("invalid json request")
		log.Printf("repositories_controller.go::apiErr.Status() %d", apiErr.Status())
		log.Printf("repositories_controller.go::apiErr.Message() %s", apiErr.Message())
		c.JSON(apiErr.Status(), apiErr)
		//log.Printf("repositories_controller.go::c.Request.Body %d",)

		return
	}
	log.Printf("repositories_controller.go var request: %s", request.Name)
	result, err := services.RepositoryService.CreateRepo(request)
	log.Printf("repositories_controller.go result: %s", result.Name)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

func CreateRepos(c *gin.Context) {
	log.Printf("repositories_controller.go Request Body: %s", c.Request.Body)
	var request []repositories.CreateRepoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		apiErr := errors.NewBadRequestError("invalid json request")
		log.Printf("repositories_controller.go::apiErr.Status() %d", apiErr.Status())
		log.Printf("repositories_controller.go::apiErr.Message() %s", apiErr.Message())
		c.JSON(apiErr.Status(), apiErr)
		//log.Printf("repositories_controller.go::c.Request.Body %d",)

		return
	}
	//log.Printf("repositories_controller.go var request: %s" , request.Name)
	result, err := services.RepositoryService.CreateRepos(request)
	//log.Printf("repositories_controller.go result: %s" , result.Name)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(result.StatusCode, result)
}
