package sql

import (
	"context"
	"git.thorn.sh/Thorn/go-vodstream/storage"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const TableName = "users"

type UserStorage struct {
	DB *sqlx.DB
}

type PaginateQueryOptions struct {
	limit  uint
	offset uint
}

func (s UserStorage) All(ctx context.Context) ([]storage.User, error) {
	var users []storage.User

	rows, err := s.DB.Queryx("SELECT * FROM %v ORDER BY id ASC", TableName)
	if err != nil {
		log.Error(err)
		return users, err
	}

	for rows.Next() {
		user := storage.User{}
		err = rows.StructScan(&user)
		if err != nil {
			log.Error(err)
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (s UserStorage) List(ctx context.Context, options PaginateQueryOptions) []storage.User {
	// STUB
	return []storage.User{}
}

func (s UserStorage) GetByID(ctx context.Context, id uint) storage.User {
	// STUB
	return storage.User{}
}

func (s UserStorage) GetByUsername(ctx context.Context, username string) storage.User {
	// STUB
	return storage.User{}
}

func (s UserStorage) GetByEmail(ctx context.Context, email string) storage.User {
	// STUB
	return storage.User{}
}

func (s UserStorage) Delete(ctx context.Context, id uint) error {
	// STUB
	return nil
}

func (s UserStorage) Insert(ctx context.Context, user storage.User) error {
	// STUB
	return nil
}

func (s UserStorage) Update(ctx context.Context, id uint, user storage.User) error {
	// STUB
	return nil
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	storage := new(UserStorage)
	storage.DB = db

	return storage
}
