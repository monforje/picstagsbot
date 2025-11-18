package repoimpl

import (
	"context"
	"errors"
	"fmt"
	"picstagsbot/internal/domain/model"
	"picstagsbot/pkg/logx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PhotoRepo struct {
	pool *pgxpool.Pool
}

func NewPhotoRepo(pool *pgxpool.Pool) *PhotoRepo {
	pr := &PhotoRepo{}

	pr.pool = pool

	return pr
}

func (r *PhotoRepo) Create(ctx context.Context, photo *model.Photo) error {
	query := `
		INSERT INTO photos (user_id, telegram_id, file_size, width, height, description, tags, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		photo.UserID,
		photo.TelegramID,
		photo.FileSize,
		photo.Width,
		photo.Height,
		photo.Description,
		photo.Tags,
		photo.CreatedAt,
	).Scan(&photo.ID)

	if err != nil {
		logx.Error("db: failed to create photo", "user_id", photo.UserID, "file_id", photo.TelegramID, "error", err)
	}
	return err
}

func (r *PhotoRepo) GetByFileID(ctx context.Context, fileID string) (*model.Photo, error) {
	query := `
		SELECT id, user_id, telegram_id, file_size, width, height, description, tags, created_at
		FROM photos
		WHERE telegram_id = $1
	`

	photo := &model.Photo{}
	err := r.pool.QueryRow(ctx, query, fileID).Scan(
		&photo.ID,
		&photo.UserID,
		&photo.TelegramID,
		&photo.FileSize,
		&photo.Width,
		&photo.Height,
		&photo.Description,
		&photo.Tags,
		&photo.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logx.Error("db: failed to get photo by file_id", "file_id", fileID, "error", err)
		return nil, err
	}

	return photo, nil
}

func (r *PhotoRepo) UpdateDescription(ctx context.Context, photoID int64, description string, tags []string) error {
	query := `
		UPDATE photos
		SET description = $1, tags = $2
		WHERE id = $3
	`

	cmd, err := r.pool.Exec(ctx, query, description, tags, photoID)
	if err != nil {
		logx.Error("db: failed to update photo description", "photo_id", photoID, "error", err)
		return err
	}

	if cmd.RowsAffected() == 0 {
		logx.Warn("db: photo not found for update", "photo_id", photoID)
		return fmt.Errorf("photo with id %d not found", photoID)
	}

	return nil
}

func (r *PhotoRepo) SearchByTag(ctx context.Context, userID int64, tag string) ([]*model.Photo, error) {
	query := `
		SELECT id, user_id, telegram_id, file_size, width, height, description, tags, created_at
		FROM photos
		WHERE user_id = $1 AND tags @> ARRAY[$2]::text[]
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID, tag)
	if err != nil {
		logx.Error("db: failed to search photos by tag", "user_id", userID, "tag", tag, "error", err)
		return nil, err
	}
	defer rows.Close()

	var photos []*model.Photo
	for rows.Next() {
		photo := &model.Photo{}
		err := rows.Scan(
			&photo.ID,
			&photo.UserID,
			&photo.TelegramID,
			&photo.FileSize,
			&photo.Width,
			&photo.Height,
			&photo.Description,
			&photo.Tags,
			&photo.CreatedAt,
		)
		if err != nil {
			logx.Error("db: failed to scan photo row", "user_id", userID, "tag", tag, "error", err)
			return nil, err
		}
		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		logx.Error("db: error iterating photo rows", "user_id", userID, "tag", tag, "error", err)
		return nil, err
	}

	return photos, nil
}
