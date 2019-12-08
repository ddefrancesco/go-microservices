package domain

import (
	"fmt"
	"go-microservices/mvc/utils"
	"log"
	"net/http"
)

var (
	users = map[uint64] *User{
		123: {Id: 123, FirstName: "Daniele", LastName: "De Francesco", Email: "ddefrancesco@gmail.com"},
	}
)

func  GetUser(userId uint64) (*User,*utils.ApplicationError) {

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
