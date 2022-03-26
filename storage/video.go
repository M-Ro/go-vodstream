package storage

import (
	"github.com/google/uuid"
	"time"
)

type Video struct {
	Id    uuid.UUID `db:"id"`
	Title string    `db:"title"`

	Length      uint64    `db:"length"` // milliseconds
	PublishedAt time.Time `db:"published_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
