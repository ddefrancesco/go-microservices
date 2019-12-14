package app

import (
	"github.com/ddefrancesco/go-microservices/mvc/controllers"
)

func mapUrls()  {
	router.GET("/users/:user_id",controllers.GetUser)
}