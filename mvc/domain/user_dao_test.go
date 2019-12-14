package domain

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetUserNotFound(t *testing.T) {
	user, err := GetUser(0)

	assert.Nil(t,user,"Non aspettiamo un utente quando usernotfound")
	assert.NotNil(t,err,"Aspettiamo un errore quando usernotfound")
	assert.EqualValues(t,http.StatusNotFound,err.StatusCode,"Aspettiamo un 404 quando usernotfound")
	assert.EqualValues(t,"user_notfound",err.Code)
	assert.EqualValues(t,"user 0 not found",err.Message)
}

func TestGetUserFound(t *testing.T) {
	user, err := GetUser(123)
	assert.NotNil(t,user,"Aspettavamo un utente quando userfound")
	assert.EqualValues(t,123,user.Id,"Aspettiamo l'utente 123 quando userfound")
	assert.Nil(t,err,"Non aspettiamo un errore quando userfound")


}