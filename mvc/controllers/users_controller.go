package controllers

import (
	"encoding/json"
	"go-microservices/mvc/services"
	"go-microservices/mvc/utils"
	"log"
	"net/http"
	"strconv"
)

func GetUser(resp http.ResponseWriter , req *http.Request)  {
	log.Print("Query get param")
	userId, err := strconv.ParseUint(req.URL.Query().Get("user_id"),10,64)
	log.Printf("user_id %v",userId)
	if err != nil {
		//return bad request
		apiErr := &utils.ApplicationError{
			Message: "user_id must be a number",
			StatusCode: http.StatusBadRequest,
			Code: "bad_request",

		}
		jsonValue,_ := json.Marshal(apiErr)
		resp.WriteHeader(apiErr.StatusCode)
		resp.Write(jsonValue)
		return
	}
	log.Printf("user_id %v",userId)
	user, apiErr := services.GetUser(userId)

	if apiErr != nil {
		//Take care of the error then return
		resp.WriteHeader(apiErr.StatusCode)
		resp.Write([]byte(apiErr.Message))
		return

	}
	log.Print("user trovato")
	jsonValue,_ := json.Marshal(user)
	resp.Write(jsonValue)
}
