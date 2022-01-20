package storage

import "time"

type User struct {
	Id         uint64
	Username   string
	Email      string
	Password   string
	PublishKey string

	CreatedAt time.Time
	UpdatedAt time.Time

	// Permissions
	CanPublish bool
	CanStream  bool
}
