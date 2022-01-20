package storage

import (
	"github.com/google/uuid"
)

type Broadcast struct {
	Id uuid.UUID
	// The user publishing this broadcast
	PublisherId uint64
	Title       string
}
