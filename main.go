package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed resources
var resources embed.FS

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "tilabt")
		w.Header().Set("Content-Type", "application/json")

		email := "stranger"
		if cookie, err := r.Cookie("email"); err == nil {
			email = cookie.Value
		}

		response := map[string]string{
			"message": "Hello " + email,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		// fmt.Fprint(w, "Hello World")

	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		okHtml, _ := resources.ReadFile("resources/ok.html")
		invalidHtml, _ := resources.ReadFile("resources/invalid.html")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		if email != "" && password != "" {
			fmt.Fprint(w, string(okHtml))
			// w.Write(okHtml)
		} else {
			fmt.Fprint(w, string(invalidHtml))
			// w.Write(invalidHtml)
		}
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
			cookie := new(http.Cookie)
			cookie.Name = "email"
			cookie.Value = email
			cookie.Path = "/"

			http.SetCookie(w, cookie)
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

	dir, _ := fs.Sub(resources, "resources")
	files := http.FileServer(http.FS(dir))
	mux.Handle("/public/", http.StripPrefix("/public", files))

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
