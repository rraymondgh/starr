// Package orbit provides functions to modify data structures among the various starr libraries.
// These functions cannot live in the starr library without causing an import cycle.
// These are wrappers around the starr library and other sub modules.
package orbit

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"golift.io/starr/lidarr"
	"golift.io/starr/prowlarr"
	"golift.io/starr/radarr"
	"golift.io/starr/readarr"
	"golift.io/starr/sonarr"
)

var ErrNotPtr = errors.New("must provide a pointer to a non-nil value")

// Copy is an easy way to copy one data structure to another.
func Copy(src interface{}, dst interface{}) error {
	if s := reflect.TypeOf(src); src == nil || s.Kind() != reflect.Ptr {
		return fmt.Errorf("copy source: %w", ErrNotPtr)
	} else if d := reflect.TypeOf(dst); dst == nil || d.Kind() != reflect.Ptr {
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

// IndexerInput represents all possible Indexer inputs.
type IndexerInput interface {
	*lidarr.IndexerInput | *prowlarr.IndexerInput | *radarr.IndexerInput |
		*readarr.IndexerInput | *sonarr.IndexerInput
}

// IndexerOutput represents all possible Indexer outputs.
type IndexerOutput interface {
	*lidarr.IndexerOutput | *prowlarr.IndexerOutput | *radarr.IndexerOutput |
		*readarr.IndexerOutput | *sonarr.IndexerOutput
}

// CopyIndexer copies an indexer from one type to another, so you may copy them among instances.
func CopyIndexer[S IndexerInput | IndexerOutput, D IndexerInput](src S, dst D, keepTags bool) (D, error) {
	if err := Copy(src, dst); err != nil {
		return dst, err
	}

	element := reflect.ValueOf(dst).Elem()
	zeroField(element.FieldByName("ID"), true)
	zeroField(element.FieldByName("Tags"), !keepTags)

	return dst, nil
}

func zeroField(field reflect.Value, really bool) {
	if really && field.CanSet() {
		field.SetZero()
	}
}
