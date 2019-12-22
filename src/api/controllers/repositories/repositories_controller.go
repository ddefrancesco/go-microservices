package repositories

import (
	"github.com/ddefrancesco/go-microservices/src/api/domain/repositories"
	"github.com/ddefrancesco/go-microservices/src/api/services"
	"github.com/ddefrancesco/go-microservices/src/api/utils/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateRepo(c *gin.Context) 	 {
	var request repositories.CreateRepoRequest
	if err := c.ShouldBindJSON(&request); err != nil{
		apiErr := errors.NewBadRequestError("invalid json request")
		c.JSON(apiErr.Status(),apiErr)
		return
	}

	result, err := services.RepositoryService.CreateRepo(request)
	if err != nil{
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
