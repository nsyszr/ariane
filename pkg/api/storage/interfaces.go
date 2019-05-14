package storage

import (
	"context"

	"github.com/nsyszr/ariane/pkg/api/runtime"
)

// Interface offers a common interface for object marshaling/unmarshaling operations and
// hides all the storage-related operations behind it.
type Interface interface {
	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
	// set to the read value from database.
	Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error
}
