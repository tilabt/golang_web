package main

import (
	"fmt"
	"net/http"
)

func main() {
	address := "localhost:8000"
	server := http.Server{
		Addr: address,
	}
	fmt.Printf("Server running on http://%s", address)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
