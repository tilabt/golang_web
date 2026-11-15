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
		query := r.URL.Query()
		name := query.Get("name")

		if name == "" {
			fmt.Fprint(w, "Please Login first")
		} else {
			fmt.Fprintf(w, "Welcome %s", name)
		}
	})

	mux.HandleFunc("/posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		fmt.Fprintln(w, "Posts with slug => ", slug)
		fmt.Fprintln(w, r.Method)
		fmt.Fprintln(w, r.RequestURI)
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
