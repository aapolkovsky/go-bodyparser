package bodyparser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type bodyKey string

// BodyKey is used to retrive body from context
var BodyKey = bodyKey("Body")

// BodyParser is a bodyparser middleware struct
type BodyParser struct {
	bodyType       reflect.Type
	onError        func(http.ResponseWriter, *http.Request)
	proceedOnError bool
}

// New creates a new bodyparser middleware for the provided type.
func New(bodyType reflect.Type) *BodyParser {
	return &BodyParser{
		bodyType:       bodyType,
		proceedOnError: false,
	}
}

// Get returns parsed body or parsing error
func Get(ctx context.Context) (interface{}, error) {
	switch v := ctx.Value(BodyKey).(type) {
	case error:
		return nil, v
	default:
		return v, nil
	}
}

// GetBody returns parsed body
func GetBody(ctx context.Context) interface{} {
	return ctx.Value(BodyKey)
}

// GetError returns parsing error
func GetError(ctx context.Context) error {
	return ctx.Value(BodyKey).(error)
}

// ProceedOnError configures bodyparser middleware to proceed on parsing error
func (bodyParser *BodyParser) ProceedOnError() *BodyParser {
	bodyParser.proceedOnError = true
	return bodyParser
}

// OnError sets an (optional) function that is called when error parsing is occured
func (bodyParser *BodyParser) OnError(onError func(http.ResponseWriter, *http.Request)) *BodyParser {
	bodyParser.onError = onError
	return bodyParser
}

// Handler apply the bodyparser middleware on the request
func (bodyParser *BodyParser) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create instance of type 'bodyType' and return pointer to that as interface
		body := reflect.New(bodyParser.bodyType).Interface()

		// Decode request body to 'body' variable
		err := json.NewDecoder(r.Body).Decode(body)
		defer r.Body.Close()

		if err != nil {
			if bodyParser.proceedOnError {
				// Populate context with error
				r = r.WithContext(context.WithValue(r.Context(), BodyKey, err))
				next.ServeHTTP(w, r)
				return
			}

			if bodyParser.onError != nil {
				// Populate context with error
				r = r.WithContext(context.WithValue(r.Context(), BodyKey, err))
				bodyParser.onError(w, r)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "400 Bad Request")
			return
		}

		// Populate context with body
		r = r.WithContext(context.WithValue(r.Context(), BodyKey, body))
		next.ServeHTTP(w, r)
	})
}
