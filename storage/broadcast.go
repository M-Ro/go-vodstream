package storage

import (
	"github.com/google/uuid"
	"time"
)

type Broadcast struct {
	Id uuid.UUID `db:"id"`

	// The user responsible for this broadcast
	BroadcasterId uint64 `db:"broadcaster_id"`

	Title string `db:"title"`

	IsActive    bool `db:"is_active"`
	IsPublished bool `db:"is_published"`

	PublishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
