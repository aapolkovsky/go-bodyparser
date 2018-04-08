package bodyparser_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/go-chi/chi"

	"github.com/aapolkovsky/go-bodyparser"
)

func ExampleNew() {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		body := bodyparser.GetBody(r.Context()).(*[]string)
		fmt.Println(*body)
	}

	// Create bodyparser middleware for slice of strings
	var parser = bodyparser.New(reflect.TypeOf([]string{}))

	router := chi.NewRouter()

	// For all routes in current router
	router.Use(parser.Handler)
	router.Post("/api", handlerFunc)

	// For specific Route
	router.With(parser.Handler).Post("/v2/api", handlerFunc)

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode([]string{"test", "strings"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://localhost/api", ioutil.NopCloser(&buffer))

	router.ServeHTTP(w, r)
	// Output: [test strings]
}

func ExampleGet_Body() {
	// Somewhere in bodyparser middleware
	ctx := context.WithValue(nil, bodyparser.BodyKey, &map[string]string{})

	// In the request handler
	i, err := bodyparser.Get(ctx)
	fmt.Println(i, err)
	// Output: &map[] <nil>
}

func ExampleGet_Error() {
	// Somewhere in bodyparser middleware
	ctx := context.WithValue(nil, bodyparser.BodyKey, errors.New("error"))

	// In the request handler
	i, err := bodyparser.Get(ctx)
	fmt.Println(i, err)
	// Output: <nil> error
}

func ExampleGetBody() {
	// Somewhere in bodyparser middleware
	ctx := context.WithValue(nil, bodyparser.BodyKey, &map[string]string{})

	// In the request handler
	result := bodyparser.GetBody(ctx).(*map[string]string)
	fmt.Println(result)
	// Output: &map[]
}

func ExampleGetError() {
	ctx := context.WithValue(nil, bodyparser.BodyKey, errors.New("error"))

	err := bodyparser.GetError(ctx)
	fmt.Println(err)
	// Output: error
}

func ExampleBodyParser_ProceedOnError_Body() {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		body, err := bodyparser.Get(r.Context())
		fmt.Println(body, err)
	}

	// Create bodyparser middleware for slice of strings with ProceedOnError option
	var parser = bodyparser.New(reflect.TypeOf([]string{})).ProceedOnError()

	router := chi.NewRouter()
	router.With(parser.Handler).Post("/api", handlerFunc)

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode([]string{"slice"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://localhost/api", ioutil.NopCloser(&buffer))

	router.ServeHTTP(w, r)
	// Output: &[slice] <nil>
}

func ExampleBodyParser_ProceedOnError_Error() {
	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		body, err := bodyparser.Get(r.Context())
		fmt.Println(body, err)
	}

	// Create bodyparser middleware for slice of strings with ProceedOnError option
	var parser = bodyparser.New(reflect.TypeOf([]string{})).ProceedOnError()

	router := chi.NewRouter()
	router.With(parser.Handler).Post("/api", handlerFunc)

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode("not slice")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://localhost/api", ioutil.NopCloser(&buffer))

	router.ServeHTTP(w, r)
	// Output: <nil> json: cannot unmarshal string into Go value of type []string
}

func ExampleBodyParser_OnError() {
	// Create bodyparser middleware for slice of strings
	var parser = bodyparser.New(reflect.TypeOf([]string{}))

	// Set error handler
	parser.OnError(func(w http.ResponseWriter, r *http.Request) {
		err := bodyparser.GetError(r.Context())
		fmt.Println(err)
	})

	router := chi.NewRouter()
	router.With(parser.Handler).Post("/api", nil)

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode("not slice")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "http://localhost/api", ioutil.NopCloser(&buffer))

	router.ServeHTTP(w, r)
	// Output: json: cannot unmarshal string into Go value of type []string
}

func ExampleBodyParser_Handler() {
	// Create bodyparser middleware for slice of strings
	var parser = bodyparser.New(reflect.TypeOf([]string{}))

	router := chi.NewRouter()

	// Use Handler method of bodyparser struct as middleware
	router.Use(parser.Handler)
}
