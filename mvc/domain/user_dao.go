package domain

import (
	"fmt"
	"github.com/ddefrancesco/go-microservices/mvc/utils"
	"log"
	"net/http"
)

var (
	users = map[uint64] *User{
		123: {Id: 123, FirstName: "Daniele", LastName: "De Francesco", Email: "ddefrancesco@gmail.com"},
	}

	UserDao usersDaoInterface
)

func init()  {
	UserDao = &userDao{}
}

type usersDaoInterface interface {
	GetUser(uint64) (*User,*utils.ApplicationError)
}
type userDao struct {

}

func (u *userDao) GetUser(userId uint64) (*User,*utils.ApplicationError) {
	log.Println("Stiamo accedendo al database")
	if user := users[userId];
	 user != nil {
		log.Printf("User: %v",user)
		return user, nil
	}

	return nil, &utils.ApplicationError{
		Message: fmt.Sprintf("user %v not found",userId),
		StatusCode: http.StatusNotFound,
		Code: "user_notfound",
	}
}
