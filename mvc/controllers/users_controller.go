package controllers

import (
	"github.com/ddefrancesco/go-microservices/mvc/services"
	"github.com/ddefrancesco/go-microservices/mvc/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func GetUser(c *gin.Context)  {
	log.Print("Query get param")
	userId, err := strconv.ParseUint(c.Param("user_id"),10,64)
	log.Printf("user_id %v",userId)
	if err != nil {
		//return bad request
		apiErr := &utils.ApplicationError{
			Message: "user_id must be a number",
			StatusCode: http.StatusBadRequest,
			Code: "bad_request",

		}
		c.JSON(apiErr.StatusCode,apiErr)
		return
	}
	log.Printf("user_id %v",userId)
	user, apiErr := services.UserService.GetUser(userId)

	if apiErr != nil {
		//Take care of the error then return
		c.JSON(apiErr.StatusCode,apiErr)
		return

	}
	log.Print("user trovato")
	c.JSON(http.StatusOK,user)
}
