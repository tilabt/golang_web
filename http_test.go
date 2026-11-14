package main_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func HelloHandler(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprint(writer, "Hello")
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
