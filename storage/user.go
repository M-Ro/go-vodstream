package storage

import (
	"github.com/M-Ro/go-vodstream/internal/domain/user"
	"time"
)

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

// UsersToDomain converts storage user models to domain models.
func UsersToDomain(users []User) []user.User {
	convertedUsers := make([]user.User, len(users))

	for i, user := range users {
		convertedUsers[i] = UserToDomain(user)
	}

	return convertedUsers
}

// UsersToStorage converts domain user models to storage models.
func UsersToStorage(users []user.User) []User {
	convertedUsers := make([]User, len(users))

	for i, user := range users {
		convertedUsers[i] = UserToStorage(user)
	}

	return convertedUsers
}

// UserToDomain converts a storage user model to a domain model.
func UserToDomain(user User) user.User {
	return user.User{
		Id:         user.Id,
		Username:   user.Username,
		Email:      user.Email,
		Password:   user.Password,
		PublishKey: user.PublishKey,
		CanPublish: user.CanPublish,
		CanStream:  user.CanStream,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

// UserToStorage converts a domain user model to a storage model.
func UserToStorage(user user.User) User {
	return User{
		Id:         user.Id,
		Username:   user.Username,
		Email:      user.Email,
		Password:   user.Password,
		PublishKey: user.PublishKey,
		CanPublish: user.CanPublish,
		CanStream:  user.CanStream,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}
