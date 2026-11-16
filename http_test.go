package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func HelloHandler(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprint(writer, "Hello")
}

func QueryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	if name == "" {
		fmt.Fprint(w, "Hello stranger")
	} else {
		fmt.Fprintf(w, "Welcome %s", name)
	}
}

func MultipleParameterValuesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	names := query["name"]

	fmt.Fprint(w, names)
	fmt.Fprint(w, strings.Join(names, " "))
}

func RequestHeader(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("content-type")
	fmt.Fprint(w, contentType)
}

func ResponseHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-Powered-By", "tilabt")

	fmt.Fprint(w, "ok")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("passowrd")

	_ = password
	fmt.Fprintf(w, "welcome back %s", email)
}

func TestHttp(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	HelloHandler(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)

	bodyString := string(body)

	assert.Equal(t, "Hello", bodyString, "should be return 'Hello'")
}

func TestQueryHttp(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?name=han&password=1234", nil)
	recorder := httptest.NewRecorder()

	QueryHandler(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	query := request.URL.Query()

	// fmt.Println(query)

	bodyString := string(body)

	assert.Contains(t, "han", query.Get("name"))
	assert.Equal(t, "Welcome han", bodyString)

}
func TestMultipleValues(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/?name=han&password=1234&name=innn", nil)
	recorder := httptest.NewRecorder()

	MultipleParameterValuesHandler(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	// query := request.URL.Query()

	bodyString := string(body)

	fmt.Println(bodyString) // [han innn]han innn

}

func TestRequestHeader(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Add("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	RequestHeader(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)
	// query := request.URL.Query()

	bodyString := string(body)

	fmt.Println(bodyString) // application/json

}

func TestResponseHeader(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	ResponseHeader(recorder, request)

	response := recorder.Result()

	responseHeader := response.Header.Get("X-Powered-By")
	fmt.Println(responseHeader)
}

func TestLogin(t *testing.T) {
	requestBody := strings.NewReader("email=user@mail&password=secret")

	request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	recorder := httptest.NewRecorder()

	LoginHandler(recorder, request)

	response := recorder.Result()
	body, _ := io.ReadAll(response.Body)

	bodyString := string(body)
	fmt.Println(bodyString)

}
