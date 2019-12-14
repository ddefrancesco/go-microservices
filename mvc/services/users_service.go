package services

import (
	"github.com/ddefrancesco/go-microservices/mvc/domain"
	"github.com/ddefrancesco/go-microservices/mvc/utils"
	"log"
)

type userService struct {

}
var (
	UserService userService
)



func (u *userService) GetUser(userId uint64) (*domain.User, *utils.ApplicationError){
	log.Printf("Nel service con %v",userId)
	user, err := domain.UserDao.GetUser(userId)
	if err != nil{
		return nil, err
	}
	return user, nil
}
