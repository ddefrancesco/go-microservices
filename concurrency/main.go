package main

import (
	"fmt"
	"github.com/ddefrancesco/go-microservices/src/api/domain/repositories"
	"os"
)

func getRequests() []repositories.CreateRepoRequest{
	requests := make([]repositories.CreateRepoRequest,0)
	return requests
}

func main()  {
	requests := getRequests()
	file , err := os.Open("/Users/ddefrancesco/Desktop/requests.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("Sto per processare %d requests",len(requests)))
}
