package user

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	"github.com/M-Ro/go-vodstream/internal/paginate"
	"github.com/M-Ro/go-vodstream/storage"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

const UsersTableName = "users"

type SqlUserStorage struct {
	DB *sqlx.DB
}

var (
	ErrNoRowsAffected = errors.New("no row found with id")
)

// insertTableName is a helper function to insert the dynamic UsersTableName property
// as bindvars cannot be used as identifiers.
func insertTableName(query string) string {
	return fmt.Sprintf(query, UsersTableName)
}

// All returns all rows in the users table.
func (s SqlUserStorage) All(ctx context.Context) ([]storage.User, error) {
	users := make([]storage.User, 0)

	rows, err := s.DB.QueryxContext(ctx, insertTableName("SELECT * FROM %s ORDER BY id ASC"))
	if err != nil {
		log.Error(err)
		return users, err
	}

	for rows.Next() {
		user := storage.User{}
		err = rows.StructScan(&user)
		if err != nil {
			log.Error(err)
			rows.Close()
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

// List returns a set of rows from the users table specified by the given pagination options.
func (s SqlUserStorage) List(ctx context.Context, options paginate.QueryOptions) ([]storage.User, error) {
	users := make([]storage.User, 0)

	sql := fmt.Sprintf(
		`SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d`,
		UsersTableName, options.Order.Field, options.Order.Method, options.Limit, options.Offset,
	)

	rows, err := s.DB.QueryxContext(ctx, sql)
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

// GetByID returns the user with the given ID, or nil on failure.
func (s SqlUserStorage) GetByID(ctx context.Context, id uint64) *storage.User {
	row := s.DB.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE id = $1`), id)

	var user storage.User
	err := row.StructScan(&user)
	if err != nil {
		if err != sql2.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &user
}

// GetByUsername returns the user with the given username, or nil on failure.
func (s SqlUserStorage) GetByUsername(ctx context.Context, username string) *storage.User {
	row := s.DB.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE username ILIKE $1`), username)

	var user storage.User
	err := row.StructScan(&user)
	if err != nil {
		if err != sql2.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &user
}

// GetByEmail returns the user with the given email, or nil on failure.
func (s SqlUserStorage) GetByEmail(ctx context.Context, email string) *storage.User {
	row := s.DB.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE email ILIKE $1`), email)

	var user storage.User
	err := row.StructScan(&user)
	if err != nil {
		if err != sql2.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &user
}

// Delete removes a user with the given ID from the table. Only returns on db error.
func (s SqlUserStorage) Delete(ctx context.Context, id uint64) error {
	_, err := s.DB.ExecContext(ctx, insertTableName(`DELETE FROM %s WHERE id = $1`), id)
	if err != nil {
		log.Errorf("SqlUserStorage::Delete: %s", err)
	}

	return err
}

// Insert takes a storage model and inserts into the db. Returns an error on failure.
// Upon insertion the ID field of the model will be set.
func (s SqlUserStorage) Insert(ctx context.Context, user *storage.User) error {
	user.CreatedAt = time.Now().Truncate(time.Microsecond)
	user.UpdatedAt = user.CreatedAt

	row := s.DB.QueryRowContext(
		ctx,
		insertTableName(`INSERT INTO %s 
			(username, email, password, publish_key, can_publish, can_stream, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`),
		user.Username, user.Email, user.Password, user.PublishKey, user.CanPublish, user.CanStream,
		user.CreatedAt, user.UpdatedAt,
	)

	err := row.Scan(&user.Id)

	return err
}

// Update takes a storage model and updates row contents for the user at the given ID.
// Returns error on failure, or if a user was not found with the given id.
func (s SqlUserStorage) Update(ctx context.Context, id uint64, user *storage.User) error {
	user.UpdatedAt = time.Now().Truncate(time.Microsecond)

	result, err := s.DB.ExecContext(
		ctx,
		insertTableName(`UPDATE %s SET 
			username=$1, email=$2, password=$3, publish_key=$4, can_publish=$5, can_stream=$6,
			created_at=$7, updated_at=$8 WHERE id=$9`),
		user.Username, user.Email, user.Password, user.PublishKey, user.CanPublish, user.CanStream, user.CreatedAt, user.UpdatedAt, id)

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
		log.Warn(ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return nil
}

// NewUserStorage instantiates a new SqlUserStorage object.
func NewUserStorage(db *sqlx.DB) *SqlUserStorage {
	newStorage := new(SqlUserStorage)
	newStorage.DB = db

	return newStorage
}
