package broadcast

import (
	"context"
	sql2 "database/sql"
	"errors"
	"fmt"
	"github.com/M-Ro/go-vodstream/internal/paginate"
	"github.com/M-Ro/go-vodstream/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"time"
)

const BroadcastsTableName = "broadcasts"

type SqlBroadcastStorage struct {
	DB *sqlx.DB
}

var (
	ErrNoRowsAffected = errors.New("no row found with id")
)

// insertTableName is a helper function to insert the dynamic BroadcastsTableName property
// as bindvars cannot be used as identifiers.
func insertTableName(query string) string {
	return fmt.Sprintf(query, BroadcastsTableName)
}

// All returns all rows in the broadcasts table.
func (s SqlBroadcastStorage) All(ctx context.Context) ([]storage.Broadcast, error) {
	broadcasts := make([]storage.Broadcast, 0)

	rows, err := s.DB.QueryxContext(ctx, insertTableName("SELECT * FROM %s ORDER BY id ASC"))
	if err != nil {
		log.Error(err)
		return broadcasts, err
	}

	for rows.Next() {
		broadcast := storage.Broadcast{}
		err = rows.StructScan(&broadcast)
		if err != nil {
			log.Error(err)
			rows.Close()
			return broadcasts, err
		}

		broadcasts = append(broadcasts, broadcast)
	}

	return broadcasts, nil
}

// List returns a set of rows from the broadcasts table specified by the given pagination options.
func (s SqlBroadcastStorage) List(ctx context.Context, options paginate.QueryOptions) ([]storage.Broadcast, error) {
	broadcasts := make([]storage.Broadcast, 0)

	sql := fmt.Sprintf(
		`SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d`,
		BroadcastsTableName, options.Order.Field, options.Order.Method, options.Limit, options.Offset,
	)

	rows, err := s.DB.QueryxContext(ctx, sql)
	if err != nil {
		log.Error(err)
		return broadcasts, err
	}

	for rows.Next() {
		broadcast := storage.Broadcast{}
		err = rows.StructScan(&broadcast)
		if err != nil {
			log.Error(err)
			return broadcasts, err
		}

		broadcasts = append(broadcasts, broadcast)
	}

	return broadcasts, nil
}

// GetByID returns the broadcast with the given ID, or nil on failure.
func (s SqlBroadcastStorage) GetByID(ctx context.Context, id uuid.UUID) *storage.Broadcast {
	row := s.DB.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE id = $1`), id)

	var broadcast storage.Broadcast
	err := row.StructScan(&broadcast)
	if err != nil {
		if err != sql2.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &broadcast
}

// Delete removes a broadcast with the given ID from the table. Only returns on db error.
func (s SqlBroadcastStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.ExecContext(ctx, insertTableName(`DELETE FROM %s WHERE id = $1`), id)
	if err != nil {
		log.Errorf("SqlBroadcastStorage::Delete: %s", err)
	}

	return err
}

// Insert takes a storage model and inserts into the db. Returns an error on failure.
// Upon insertion the ID field of the model will be set.
func (s SqlBroadcastStorage) Insert(ctx context.Context, broadcast *storage.Broadcast) error {
	broadcast.CreatedAt = time.Now().Truncate(time.Microsecond)
	broadcast.UpdatedAt = broadcast.CreatedAt

	row := s.DB.QueryRowContext(
		ctx,
		insertTableName(`INSERT INTO %s 
			(broadcaster_id, title, is_active, is_published, published_at, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id`),
		broadcast.BroadcasterId, broadcast.Title, broadcast.IsActive, broadcast.IsPublished, broadcast.PublishedAt,
		broadcast.CreatedAt, broadcast.UpdatedAt,
	)

	err := row.Scan(&broadcast.Id)

	return err
}

// Update takes a storage model and updates row contents for the user at the given ID.
// Returns error on failure, or if a user was not found with the given id.
func (s SqlBroadcastStorage) Update(ctx context.Context, id uuid.UUID, broadcast *storage.Broadcast) error {
	broadcast.UpdatedAt = time.Now().Truncate(time.Microsecond)

	result, err := s.DB.ExecContext(
		ctx,
		insertTableName(`UPDATE %s SET 
			broadcaster_id=$1, title=$2, is_active=$3, is_published=$4, published_at=$5, created_at=$6,
			updated_at=$7, WHERE id=$8`),
		broadcast.BroadcasterId, broadcast.Title, broadcast.IsActive, broadcast.IsPublished, broadcast.PublishedAt,
		broadcast.CreatedAt, broadcast.UpdatedAt, id)

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

// NewBroadcastStorage instantiates a new SqlBroadcastStorage object.
func NewBroadcastStorage(db *sqlx.DB) *SqlBroadcastStorage {
	newStorage := new(SqlBroadcastStorage)
	newStorage.DB = db

	return newStorage
}
