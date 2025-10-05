package mysql

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	entropyOnce sync.Once
	entropy     *ulid.MonotonicEntropy
)

// NewID generates a monotonic ULID suitable for use as a primary key.
func NewID() string {
	entropyOnce.Do(func() {
		entropy = ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	})

	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
