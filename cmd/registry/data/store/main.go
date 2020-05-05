package store

import (
	"errors"

	"github.com/noptics/services/cmd/registry/data"
	"github.com/noptics/services/cmd/registry/data/dynamo"
)

// NewStore returns a configured data store backend.
func New(storeType string, config map[string]string) (data.Store, error) {
	switch storeType {
	case "dynamo":
		return dynamo.New(config)
	default:
		return nil, errors.New("Unknown data store provider type")
	}
}
