package broadcast_vod

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

const BroadcastVodsTableName = "broadcast_vods"

type SqlBroadcastVodStorage struct {
	DB *sqlx.DB
}

var (
	ErrNoRowsAffected = errors.New("no row found with id")
)

// insertTableName is a helper function to insert the dynamic BroadcastVodsTableName property
// as bindvars cannot be used as identifiers.
func insertTableName(query string) string {
	return fmt.Sprintf(query, BroadcastVodsTableName)
}

// All returns all rows in the broadcast_vods table.
func (s SqlBroadcastVodStorage) All(ctx context.Context) ([]storage.BroadcastVod, error) {
	broadcastVods := make([]storage.BroadcastVod, 0)

	rows, err := s.DB.QueryxContext(ctx, insertTableName("SELECT * FROM %s ORDER BY id ASC"))
	if err != nil {
		log.Error(err)
		return broadcastVods, err
	}

	for rows.Next() {
		broadcastVod := storage.BroadcastVod{}
		err = rows.StructScan(&broadcastVod)
		if err != nil {
			log.Error(err)
			rows.Close()
			return broadcastVods, err
		}

		broadcastVods = append(broadcastVods, broadcastVod)
	}

	return broadcastVods, nil
}

// List returns a set of rows from the broadcast_vods table specified by the given pagination options.
func (s SqlBroadcastVodStorage) List(ctx context.Context, options paginate.QueryOptions) ([]storage.BroadcastVod, error) {
	broadcastVods := make([]storage.BroadcastVod, 0)

	sql := fmt.Sprintf(
		`SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d`,
		BroadcastVodsTableName, options.Order.Field, options.Order.Method, options.Limit, options.Offset,
	)

	rows, err := s.DB.QueryxContext(ctx, sql)
	if err != nil {
		log.Error(err)
		return broadcastVods, err
	}

	for rows.Next() {
		broadcastVod := storage.BroadcastVod{}
		err = rows.StructScan(&broadcastVod)
		if err != nil {
			log.Error(err)
			return broadcastVods, err
		}

		broadcastVods = append(broadcastVods, broadcastVod)
	}

	return broadcastVods, nil
}

// GetByID returns the broadcastVod with the given ID, or nil on failure.
func (s SqlBroadcastVodStorage) GetByID(ctx context.Context, id uuid.UUID) *storage.BroadcastVod {
	row := s.DB.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE id = $1`), id)

	var broadcastVod storage.BroadcastVod
	err := row.StructScan(&broadcastVod)
	if err != nil {
		if err != sql2.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &broadcastVod
}

// Delete removes a broadcastVod with the given ID from the table. Only returns on db error.
func (s SqlBroadcastVodStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.ExecContext(ctx, insertTableName(`DELETE FROM %s WHERE id = $1`), id)
	if err != nil {
		log.Errorf("SqlBroadcastVodStorage::Delete: %s", err)
	}

	return err
}

// Insert takes a storage model and inserts into the db. Returns an error on failure.
// Upon insertion the ID field of the model will be set.
func (s SqlBroadcastVodStorage) Insert(ctx context.Context, broadcastVod *storage.BroadcastVod) error {
	broadcastVod.CreatedAt = time.Now().Truncate(time.Microsecond)
	broadcastVod.UpdatedAt = broadcastVod.CreatedAt

	row := s.DB.QueryRowContext(
		ctx,
		insertTableName(`INSERT INTO %s 
			(stream_id, video_id, published_at, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5) RETURNING id`),
		broadcastVod.StreamId, broadcastVod.VideoId, broadcastVod.PublishedAt,
		broadcastVod.CreatedAt, broadcastVod.UpdatedAt,
	)

	err := row.Scan(&broadcastVod.Id)

	return err
}

// Update takes a storage model and updates row contents for the user at the given ID.
// Returns error on failure, or if a user was not found with the given id.
func (s SqlBroadcastVodStorage) Update(ctx context.Context, id uuid.UUID, broadcastVod *storage.BroadcastVod) error {
	broadcastVod.UpdatedAt = time.Now().Truncate(time.Microsecond)

	result, err := s.DB.ExecContext(
		ctx,
		insertTableName(`UPDATE %s SET 
			stream_id=$1, video_id=$2, published_at=$3, created_at=$4, created_at=$5, WHERE id=$6`),
		broadcastVod.StreamId, broadcastVod.VideoId, broadcastVod.PublishedAt,
		broadcastVod.CreatedAt, broadcastVod.UpdatedAt, id)

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

// NewBroadcastVodStorage instantiates a new SqlBroadcastVodStorage object.
func NewBroadcastVodStorage(db *sqlx.DB) *SqlBroadcastVodStorage {
	newStorage := new(SqlBroadcastVodStorage)
	newStorage.DB = db

	return newStorage
}
