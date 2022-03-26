package user

import "time"

type User struct {
	Id         uint64
	Username   string
	Email      string
	Password   string
	PublishKey string

	// Permissions
	CanPublish bool
	CanStream  bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
