package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "tilabt")
		w.Header().Set("Content-Type", "application/json")

		response := map[string]string{
			"message": "Hello World",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// fmt.Fprint(w, "Hello World")

	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		//POST /login HTTP/1.1
		// Content-Length: 42
		// Content-Type: application/x-www-form-urlencoded
		// Host: localhost:8000
		// User-Agent: HTTPie

		// email=hanasa%40gmail.com&password=password

		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		w.Header().Set("Content-type", "application/json")
		///opt/homebrew/Cellar/go/1.27.3/libexec/src/net/http/status.go
		// r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")
		response := map[string]any{}

		if email != "" && password != "" {
			w.WriteHeader(http.StatusOK)
			response = map[string]any{
				"message": fmt.Sprintf("Welcome back %s", email),
				"data": map[string]string{
					"email":    email,
					"password": password,
				},
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			response = map[string]any{
				"message": "Email and Password are required",
			}
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			panic(err)
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
