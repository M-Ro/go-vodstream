package video

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

const VideosTableName = "videos"

type SqlVideoStorage struct {
	DB *sqlx.DB
}

var (
	ErrNoRowsAffected = errors.New("no row found with id")
)

// insertTableName is a helper function to insert the dynamic VideosTableName property
// as bindvars cannot be used as identifiers.
func insertTableName(query string) string {
	return fmt.Sprintf(query, VideosTableName)
}

// All returns all rows in the videos table.
func (s SqlVideoStorage) All(ctx context.Context) ([]storage.Video, error) {
	videos := make([]storage.Video, 0)

	rows, err := s.DB.QueryxContext(ctx, insertTableName("SELECT * FROM %s ORDER BY id ASC"))
	if err != nil {
		log.Error(err)
		return videos, err
	}

	for rows.Next() {
		video := storage.Video{}
		err = rows.StructScan(&video)
		if err != nil {
			log.Error(err)
			rows.Close()
			return videos, err
		}

		videos = append(videos, video)
	}

	return videos, nil
}

// List returns a set of rows from the broadcast_vods table specified by the given pagination options.
func (s SqlVideoStorage) List(ctx context.Context, options paginate.QueryOptions) ([]storage.Video, error) {
	videos := make([]storage.Video, 0)

	sql := fmt.Sprintf(
		`SELECT * FROM %s ORDER BY %s %s LIMIT %d OFFSET %d`,
		VideosTableName, options.Order.Field, options.Order.Method, options.Limit, options.Offset,
	)

	rows, err := s.DB.QueryxContext(ctx, sql)
	if err != nil {
		log.Error(err)
		return videos, err
	}

	for rows.Next() {
		video := storage.Video{}
		err = rows.StructScan(&video)
		if err != nil {
			log.Error(err)
			return videos, err
		}

		videos = append(videos, video)
	}

	return videos, nil
}

// GetByID returns the broadcastVod with the given ID, or nil on failure.
func (s SqlVideoStorage) GetByID(ctx context.Context, id uuid.UUID) *storage.Video {
	row := s.DB.QueryRowxContext(ctx, insertTableName(`SELECT * from %s WHERE id = $1`), id)

	var video storage.Video
	err := row.StructScan(&video)
	if err != nil {
		if err != sql2.ErrNoRows {
			log.Error(err)
		}
		return nil
	}

	return &video
}

// Delete removes a broadcastVod with the given ID from the table. Only returns on db error.
func (s SqlVideoStorage) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.DB.ExecContext(ctx, insertTableName(`DELETE FROM %s WHERE id = $1`), id)
	if err != nil {
		log.Errorf("SqlVideoStorage::Delete: %s", err)
	}

	return err
}

// Insert takes a storage model and inserts into the db. Returns an error on failure.
// Upon insertion the ID field of the model will be set.
func (s SqlVideoStorage) Insert(ctx context.Context, video *storage.Video) error {
	video.CreatedAt = time.Now().Truncate(time.Microsecond)
	video.UpdatedAt = video.CreatedAt

	row := s.DB.QueryRowContext(
		ctx,
		insertTableName(`INSERT INTO %s 
			(title, length, published_at, created_at, updated_at)
			VALUES ($1,$2,$3,$4,$5) RETURNING id`),
		video.Title, video.Length, video.PublishedAt,
		video.CreatedAt, video.UpdatedAt,
	)

	err := row.Scan(&video.Id)

	return err
}

// Update takes a storage model and updates row contents for the user at the given ID.
// Returns error on failure, or if a user was not found with the given id.
func (s SqlVideoStorage) Update(ctx context.Context, id uuid.UUID, video *storage.Video) error {
	video.UpdatedAt = time.Now().Truncate(time.Microsecond)

	result, err := s.DB.ExecContext(
		ctx,
		insertTableName(`UPDATE %s SET 
			title=$1, length=$2, published_at=$3, created_at=$4, created_at=$5, WHERE id=$6`),
		video.Title, video.Length, video.PublishedAt,
		video.CreatedAt, video.UpdatedAt, id)

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

// NewVideoStorage instantiates a new SqlVideoStorage object.
func NewVideoStorage(db *sqlx.DB) *SqlVideoStorage {
	newStorage := new(SqlVideoStorage)
	newStorage.DB = db

	return newStorage
}
