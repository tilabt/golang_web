package main

import (
	"fmt"
	"net/http"
)

func main() {
	var hanlder http.HandlerFunc = func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "Hello world")
	}
	address := "localhost:8000"

	server := http.Server{
		Addr:    address,
		Handler: hanlder,
	}

	fmt.Printf("Server running on http://%s", address)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
