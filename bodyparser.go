package bodyparser

import (
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
	if !isValidType(bodyType) {
		return nil
	}

	return &BodyParser{
		bodyType:       bodyType,
		proceedOnError: false,
	}
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
		body, err := Parse(r.Body, bodyParser.bodyType)

		if err != nil {
			if bodyParser.proceedOnError {
				// Populate context with error
				r = r.WithContext(Put(r.Context(), err))

				// Serve next
				next.ServeHTTP(w, r)
				return
			}

			if bodyParser.onError != nil {
				// Populate context with error
				r = r.WithContext(Put(r.Context(), err))

				// Call special error handler
				bodyParser.onError(w, r)
				return
			}

			// Respond with "Bad Request" error
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Close request body
		defer r.Body.Close()

		// Populate context with body
		r = r.WithContext(Put(r.Context(), body))

		// Serve next
		next.ServeHTTP(w, r)
	})
}
