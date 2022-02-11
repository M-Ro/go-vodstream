package sql

import (
	"context"
	"errors"
	"fmt"
	"git.thorn.sh/Thorn/go-vodstream/internal/domain"
	"git.thorn.sh/Thorn/go-vodstream/storage"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

const TableName = "users"

type UserStorage struct {
	DB *sqlx.DB
}

var (
	ErrNoRowsAffected = errors.New("no row found with id")
)

// insertTableName is a helper function to insert the dynamic TableName property
// as bindvars cannot be used as identifiers.
func insertTableName(query string) string {
	return fmt.Sprintf(query, TableName)
}

// All returns all rows in the users table.
func (s UserStorage) All(ctx context.Context) ([]storage.User, error) {
	var users []storage.User

	rows, err := s.DB.QueryContext(ctx, insertTableName("SELECT * FROM %s ORDER BY id ASC"))
	if err != nil {
		log.Error(err)
		return users, err
	}

	for rows.Next() {
		user := storage.User{}
		err = rows.Scan(&user)
		if err != nil {
			log.Error(err)
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

// List returns a set of rows from the users table specified by the given pagination options.
func (s UserStorage) List(ctx context.Context, options domain.PaginateQueryOptions) ([]storage.User, error) {
	var users []storage.User

	sql := fmt.Sprintf(
		`SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d`,
		TableName, options.Order.Field, options.Order.Method, options.Limit, options.Offset,
	)

	rows, err := s.DB.QueryContext(ctx, sql)
	if err != nil {
		log.Error(err)
		return users, err
	}

	for rows.Next() {
		user := storage.User{}
		err = rows.Scan(&user)
		if err != nil {
			log.Error(err)
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetByID returns the user with the given ID, or nil on failure.
func (s UserStorage) GetByID(ctx context.Context, id uint) *storage.User {
	row := db.QueryRowContext(ctx, insertTableName(`SELECT * from %s WHERE id = $1`), id)

	var user storage.User
	err := row.Scan(&user)
	if err != nil {
		log.Error(err)
		return nil
	}

	return &user
}

// GetByUsername returns the user with the given username, or nil on failure.
func (s UserStorage) GetByUsername(ctx context.Context, username string) *storage.User {
	row := db.QueryRowContext(ctx, insertTableName(`SELECT * from %s WHERE username = $1`), username)

	var user storage.User
	err := row.Scan(&user)
	if err != nil {
		log.Error(err)
		return nil
	}

	return &user
}

// GetByEmail returns the user with the given email, or nil on failure.
func (s UserStorage) GetByEmail(ctx context.Context, email string) *storage.User {
	row := db.QueryRowContext(ctx, insertTableName(`SELECT * from %s WHERE email = $1`), email)

	var user storage.User
	err := row.Scan(&user)
	if err != nil {
		log.Error(err)
		return nil
	}

	return &user
}

// Delete removes a user with the given ID from the table. Only returns on db error.
func (s UserStorage) Delete(ctx context.Context, id uint) error {
	_, err := s.DB.ExecContext(ctx, insertTableName(`DELETE FROM %s WHERE id = $1`), id)
	if err != nil {
		log.Errorf("UserStorage::Delete: %s", err)
	}

	return err
}

// Insert takes a storage model and inserts into the db. Returns an error on failure.
// Upon insertion the ID field of the model will be set.
func (s UserStorage) Insert(ctx context.Context, user *storage.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	row := s.DB.QueryRowContext(
		ctx,
		insertTableName(`INSERT INTO %s 
			(username, email, password, publish_key, can_publish, can_stream, created_at, updated_at)
			VALUES ($1 $2 $3 $4 $5 $6 $7 $8)`),
		user.Username, user.Email, user.Password, user.PublishKey, user.CanPublish, user.CanStream,
		user.CreatedAt, user.UpdatedAt,
	)

	err := row.Scan(&user.Id)

	return err
}

// Update takes a storage model and updates row contents for the user at the given ID.
// Returns error on failure, or if a user was not found with the given id.
func (s UserStorage) Update(ctx context.Context, id uint, user *storage.User) error {
	user.UpdatedAt = time.Now()

	result, err := s.DB.ExecContext(
		ctx,
		insertTableName(`UPDATE %s WHERE id=$1 SET 
			username=$2, email=$3, password=$4, publish_key=$5, can_publish=$6, can_stream=$7,
			created_at=$8, updated_at=$9`),
		id, user.Username, user.Email, user.Password, user.PublishKey, user.CanPublish, user.CanStream, user.CreatedAt, user.UpdatedAt)

	if err != nil {
		log.Error(err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}

	if rows != 1 {
		log.Error(ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return nil
}

// NewUserStorage instantiates a new UserStorage object.
func NewUserStorage(db *sqlx.DB) *UserStorage {
	storage := new(UserStorage)
	storage.DB = db

	return storage
}
