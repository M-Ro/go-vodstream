package storage

import "time"

type User struct {
	Id         uint64 `db:"id"`
	Username   string `db:"username"`
	Email      string `db:"email"`
	Password   string `db:"password"`
	PublishKey string `db:"publish_key"`

	// Permissions
	CanPublish bool `db:"can_publish"`
	CanStream  bool `db:"can_stream"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
