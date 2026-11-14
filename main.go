package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Login")
	})

	mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Posts")
	})
	mux.HandleFunc("/posts/likes/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Posts Likes")
	})

	address := "localhost:8000"

	server := http.Server{
		Addr:    address,
		Handler: mux,
	}

	fmt.Printf("Server running on http://%s", address)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
