package starr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var ErrNotPtr = errors.New("must provide a pointer to a struct")

// Copy is an easy way to copy one data structure to another.
func Copy(src interface{}, dst interface{}) error {
	if s := reflect.TypeOf(src); s.Kind() != reflect.Ptr ||
		s.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("copy source: %w", ErrNotPtr)
	} else if d := reflect.TypeOf(dst); d.Kind() != reflect.Ptr ||
		d.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("copy destination: %w", ErrNotPtr)
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(src); err != nil {
		return fmt.Errorf("encoding: %w", err)
	}

	if err := json.NewDecoder(&buf).Decode(dst); err != nil {
		return fmt.Errorf("decoding: %w", err)
	}

	return nil
}
