package user

import (
	"context"
	"errors"
	"git.thorn.sh/Thorn/go-vodstream/internal/domain"
	"git.thorn.sh/Thorn/go-vodstream/internal/paginate"
	"git.thorn.sh/Thorn/go-vodstream/storage"
)

var (
	ErrUserNotFound = errors.New("no user found")
)

type StorageProvider interface {
	All(ctx context.Context) ([]storage.User, error)
	List(ctx context.Context, options paginate.QueryOptions) ([]storage.User, error)
	GetByID(ctx context.Context, id uint64) *storage.User
	GetByUsername(ctx context.Context, username string) *storage.User
	GetByEmail(ctx context.Context, email string) *storage.User
	Delete(ctx context.Context, id uint64) error
	Insert(ctx context.Context, user *storage.User) error
	Update(ctx context.Context, id uint64, user *storage.User) error
}

type Repository struct {
	StorageProvider StorageProvider
}

// All returns all users in the repository.
func (r Repository) All(ctx context.Context) ([]domain.User, error) {
	users, err := r.StorageProvider.All(ctx)
	if err != nil {
		return []domain.User{}, err
	}

	return storage.UsersToDomain(users), nil
}

// List returns a set of users specified by the provided QueryOptions.
func (r Repository) List(ctx context.Context, options paginate.QueryOptions) ([]domain.User, error) {
	users, err := r.StorageProvider.List(ctx, options)
	if err != nil {
		return []domain.User{}, err
	}

	return storage.UsersToDomain(users), nil
}

// GetByID returns the user with the given ID, or returns an error.
func (r Repository) GetByID(ctx context.Context, id uint64) (domain.User, error) {
	user := r.StorageProvider.GetByID(ctx, id)
	if user == nil {
		return domain.User{}, ErrUserNotFound
	}

	return storage.UserToDomain(*user), nil
}

// GetByUsername returns the user with the given Username, or returns an error.
func (r Repository) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	user := r.StorageProvider.GetByUsername(ctx, username)
	if user == nil {
		return domain.User{}, ErrUserNotFound
	}

	return storage.UserToDomain(*user), nil
}

// GetByEmail returns the user with the given Email, or returns an error.
func (r Repository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	user := r.StorageProvider.GetByEmail(ctx, email)
	if user == nil {
		return domain.User{}, ErrUserNotFound
	}

	return storage.UserToDomain(*user), nil
}

// Delete removes a user with the given ID from the table. Returns an error on failure.
func (r Repository) Delete(ctx context.Context, id uint64) error {
	return r.StorageProvider.Delete(ctx, id)
}

// Insert takes a domain model and inserts it to the storage provider.
// After successful insertion the user ID field should be filled, alternatively an error is returned.
func (r Repository) Insert(ctx context.Context, user *domain.User) error {
	storageUser := storage.UserToStorage(*user)

	err := r.StorageProvider.Insert(ctx, &storageUser)
	if err != nil {
		return err
	}

	user.Id = storageUser.Id

	return nil
}

// Update takes a user and updates the record within the StorageProvider.
func (r Repository) Update(ctx context.Context, id uint64, user domain.User) (domain.User, error) {
	storageUser := storage.UserToStorage(user)

	err := r.StorageProvider.Update(ctx, id, &storageUser)
	if err != nil {
		return domain.User{}, err
	}

	return storage.UserToDomain(storageUser), nil
}

func NewRepository(s StorageProvider) Repository {
	return Repository{
		StorageProvider: s,
	}
}
