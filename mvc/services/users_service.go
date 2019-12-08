package services

import (
	"go-microservices/mvc/domain"
	"go-microservices/mvc/utils"
	"log"
)

func GetUser(userId uint64) (*domain.User, *utils.ApplicationError){
	log.Printf("Nel service con %v",userId)
	return  domain.GetUser(userId)
}
