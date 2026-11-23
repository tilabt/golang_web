package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed resources
var resources embed.FS

//go:embed templates
var templates embed.FS

type ProfileAddress struct {
	Country, Street string
	No              int16
}
type Profile struct {
	Title, Name string
	Age         int16
	Address     *ProfileAddress
}

type User struct {
	Title, Name string
}

func (u User) SayHello() string {
	return fmt.Sprintf("Hello %s from FN", u.Name)
}

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

	mux.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		templateText := "<html><body>{{.}}</body></html>"
		t := template.Must(template.New("ME").Parse(templateText))

		if err := t.ExecuteTemplate(w, "ME", "Welcome back"); err != nil {
			log.Fatalf("Failed execute template: %v", err)
		}
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("./templates/hello.gohtml"))

		if err := t.ExecuteTemplate(w, "hello.gohtml", "Hello cuy"); err != nil {
			log.Fatalf("Failed execute template: %v", err)
		}
	})

	mux.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		// using template DIR instead of file
		t := template.Must(template.ParseGlob("./templates/*.gohtml"))

		if err := t.ExecuteTemplate(w, "hello.gohtml", "Hi cuy"); err != nil {
			log.Fatalf("Failed execute template: %v", err)
		}
	})

	mux.HandleFunc("/yuhu", func(w http.ResponseWriter, r *http.Request) {
		// embed file
		t := template.Must(template.ParseFS(templates, "templates/*.gohtml"))

		if err := t.ExecuteTemplate(w, "hello.gohtml", "Yuhu cuy"); err != nil {
			log.Fatalf("Failed execute template: %v", err)
		}
	})

	// using data map
	mux.HandleFunc("/me", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("./templates/me.gohtml"))
		data := map[string]any{
			"Title": "Me Profile",
			"Name":  "John Doe",
			"Age":   28,
			"Address": map[string]any{
				"Country": "NoWhere",
				"Street":  "Sreeeet",
				"No":      1234,
			},
		}
		if err := t.ExecuteTemplate(w, "me.gohtml", data); err != nil {
			log.Fatalf("Failed execute template: %v", err)
		}

	})

	// using data struct
	mux.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("./templates/me.gohtml"))
		data := Profile{
			Title: "My Profile",
			Name:  "John Doe",
			Age:   28,
			// Address: ProfileAddress{
			// 	Country: "NoWhere",
			// 	Street:  "Sreeeet",
			// 	No:      1234,
			// },
		}
		if err := t.ExecuteTemplate(w, "me.gohtml", data); err != nil {
			log.Fatalf("Failed execute template: %v", err)
		}

	})

	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		fn := template.FuncMap{
			"upper": func(value string) string {
				return strings.ToUpper(value)
			},
			// suffix first then value
			"append": func(suffix, value string) string {
				return value + suffix
			},
			// prefix first then value
			"prepend": func(prefix, value string) string {
				return prefix + value
			},
		}

		t := template.Must(template.
			New("base").
			Funcs(fn).
			ParseFiles(
				"./partials/layouts/base.gohtml",
				"./templates/user.gohtml",
				"./partials/nav.gohtml",
			))

		t.ExecuteTemplate(w, "base", User{
			Title: "My User",
			Name:  name,
		})
	})

	mux.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		rate := 2.1

		t := template.Must(template.ParseFiles("./templates/rate.gohtml"))
		if err := t.ExecuteTemplate(w, "rate.gohtml", map[string]any{
			"Rate": rate,
		}); err != nil {
			log.Printf("render rate template: %v", err)
			http.Error(w, "Gagal merender halaman", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("/posts", func(w http.ResponseWriter, r *http.Request) {
		t := template.Must(template.ParseFiles("./templates/posts.gohtml"))

		t.ExecuteTemplate(w, "posts.gohtml", map[string]any{
			// "Posts": []string{
			// 	"Post-1", "Post-2", "Post-3",
			// },
		})
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

	address := "localhost:3000"

	server := http.Server{
		Addr:    address,
		Handler: mux,
	}

	fmt.Printf("Server running on http://%s", address)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}

}
