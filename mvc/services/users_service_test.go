package services

import (
	"github.com/stretchr/testify/assert"
	"github.com/ddefrancesco/go-microservices/mvc/domain"
	"github.com/ddefrancesco/go-microservices/mvc/utils"
	"net/http"
	"testing"
)

var (
	userDaoMock usersDaoMock
	getUserFunction func(userId uint64) (*domain.User,*utils.ApplicationError)
)
func init()  {
	domain.UserDao = &usersDaoMock{}
}

type usersDaoMock struct {
}

func (m *usersDaoMock)GetUser(userId uint64) (*domain.User,*utils.ApplicationError){
	return getUserFunction(userId)
}
func TestUserNotFoundInDatabase(t *testing.T)  {
	getUserFunction = func(userId uint64) (*domain.User,*utils.ApplicationError){
		return nil,&utils.ApplicationError{
			StatusCode:http.StatusNotFound,
			Message:"user 0 not found",
		}
	}
	user, err := UserService.GetUser(0)
	assert.Nil(t, user)
	assert.NotNil(t,err)
	assert.EqualValues(t, http.StatusNotFound,err.StatusCode)
	assert.EqualValues(t, "user 0 not found", err.Message)

}
func TestUserNoError(t *testing.T)  {
	getUserFunction = func(userId uint64) (*domain.User,*utils.ApplicationError){
		return &domain.User{
			Id:123,
		},nil
	}
	user, err := UserService.GetUser(123)
	assert.Nil(t, err)
	assert.NotNil(t,user)
	assert.EqualValues(t, 123, user.Id)
}