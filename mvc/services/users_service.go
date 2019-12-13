package services

import (
	"github.com/ddefrancesco/go-microservices/mvc/domain"
	"github.com/ddefrancesco/go-microservices/mvc/utils"
	"log"
)

func GetUser(userId uint64) (*domain.User, *utils.ApplicationError){
	log.Printf("Nel service con %v",userId)
	return  domain.GetUser(userId)
}
