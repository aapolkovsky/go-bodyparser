package bodyparser

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"reflect"
)

var invalidTypes = []reflect.Kind{
	reflect.Invalid,
	reflect.Complex64,
	reflect.Complex128,
	reflect.Chan,
	reflect.Func,
	reflect.Ptr,
	reflect.UnsafePointer,
}

// ErrInvalidType indicates that bodyType is invalid
var ErrInvalidType = errors.New("bodyparser: Invalid type of body")

func isValidType(bodyType reflect.Type) bool {
	kind := bodyType.Kind()

	for _, t := range invalidTypes {
		if kind == t {
			return false
		}
	}
	return true
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

// Put puts data (body or error) to context
func Put(ctx context.Context, data interface{}) context.Context {
	return context.WithValue(ctx, BodyKey, data)
}

// NewOfType creates new instance of bodyType and returns it as interface
func NewOfType(bodyType reflect.Type) (interface{}, error) {
	if !isValidType(bodyType) {
		return nil, ErrInvalidType
	}

	// Put pointer to body
	return reflect.New(bodyType).Interface(), nil
}

// Parse parses bytes from reader to new instance of bodyType
func Parse(reader io.Reader, bodyType reflect.Type) (interface{}, error) {
	body, err := NewOfType(bodyType)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(reader).Decode(body)

	return body, err
}
