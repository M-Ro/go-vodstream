package storage

import (
	"github.com/google/uuid"
	"time"
)

type BroadcastVod struct {
	Id uuid.UUID `db:"id"`

	StreamId uuid.UUID `db:"stream_id"`
	VideoId  uuid.UUID `db:"video_id"`

	PublishedAt time.Time `db:"published_at"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
