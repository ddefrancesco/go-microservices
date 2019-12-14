package services

import (
	"github.com/ddefrancesco/go-microservices/mvc/domain"
	"github.com/ddefrancesco/go-microservices/mvc/utils"
	"log"
	"net/http"
)

func GetItem(itemId uint64) (*domain.Item, *utils.ApplicationError){
	log.Printf("Nel service con %v",itemId)
	return  nil, &utils.ApplicationError{
		Message:    "implement me",
		StatusCode: http.StatusInternalServerError,
		Code:       "",
	}
}
