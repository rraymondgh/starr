package galaxy

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

type IndexerInput interface {
	*lidarr.IndexerInput | *prowlarr.IndexerInput | *radarr.IndexerInput |
		*readarr.IndexerInput | *sonarr.IndexerInput
}

type IndexerOutput interface {
	*lidarr.IndexerOutput | *prowlarr.IndexerOutput | *radarr.IndexerOutput |
		*readarr.IndexerOutput | *sonarr.IndexerOutput
}

// CopyIndexer copies an indexer from one type to another, so you may copy them among instances.
func CopyIndexer[S IndexerInput | IndexerOutput, D IndexerInput](src S, dst D, keepTags bool) (D, error) {
	if err := Copy(src, dst); err != nil {
		return dst, err
	}

	dstElem := reflect.ValueOf(dst).Elem()
	dstElem.FieldByName("ID").SetZero()

	if !keepTags {
		dstElem.FieldByName("Tags").SetZero()
	}

	return dst, nil
}

// Must can be used to avoid checking an error.
func Must[S any](input S, err error) S {
	if err != nil {
		panic("Must failed: " + err.Error())
	}

	return input
}
