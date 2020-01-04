package repositories

import (
	"github.com/ddefrancesco/go-microservices/src/api/client/restclient"
	"github.com/ddefrancesco/go-microservices/src/api/utils/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M){
	restclient.StartMockups()
	os.Exit(m.Run())
}

func TestCreateRepoInvalidJsonBody(t *testing.T) {

	response := httptest.NewRecorder()
	c,_ := gin.CreateTestContext(response)

	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(``))

	c.Request = request

	CreateRepo(c)


	assert.EqualValues(t, http.StatusBadRequest, response.Code)
	assert.NotEmpty(t,response.Body)
	log.Printf("repositories_controller_test.go: response body: %d" ,response.Body.Len())
	apiErr, err := errors.NewApiErrorFromBytes(response.Body.Bytes())
	assert.Nil(t,err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t,http.StatusBadRequest, apiErr.Status())
	//assert.EqualValues(t,"",apiErr.Message())
	log.Printf("repositories_controller_test.go: apiErr: Status %d, Message %s ", apiErr.Status(),apiErr.Message())
}

func TestCreateRepoErrorFromGithub(t *testing.T) {

	response := httptest.NewRecorder()
	c,_ := gin.CreateTestContext(response)

	request, _ := http.NewRequest(http.MethodPost, "/repositories", strings.NewReader(`{"name": "testing"}`))

	c.Request = request

	restclient.FlushMockups()

	restclient.AddMockup(&restclient.Mock{
		Url:        "https://api.github.com/user/repos",
		HttpMethod: http.MethodPost,
		Response: &http.Response{
			StatusCode:       http.StatusUnauthorized,
			Body:             ioutil.NopCloser( strings.NewReader(`{"message": "Requires authentication", "documentation_url": "https://developer.github.com/v3/repos/#create"}`)),
		},
		Err: nil,
	})



	CreateRepo(c)


	assert.EqualValues(t, http.StatusUnauthorized, response.Code)
	assert.NotEmpty(t,response.Body)
	log.Printf("repositories_controller_test.go: response body: %d" ,response.Body.Len())
	apiErr, err := errors.NewApiErrorFromBytes(response.Body.Bytes())
	assert.Nil(t,err)
	assert.NotNil(t, apiErr)
	assert.EqualValues(t,http.StatusUnauthorized, apiErr.Status())
	//assert.EqualValues(t,"",apiErr.Message())
	log.Printf("repositories_controller_test.go: apiErr: Status %d, Message %s ", apiErr.Status(),apiErr.Message())
}