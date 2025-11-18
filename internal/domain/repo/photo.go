package repo

import (
	"context"
	"picstagsbot/internal/domain/model"
)

type PhotoRepo interface {
	Create(ctx context.Context, photo *model.Photo) error
	GetByFileID(ctx context.Context, fileID string) (*model.Photo, error)
	UpdateDescription(ctx context.Context, photoID int64, description string, tags []string) error
	SearchByTag(ctx context.Context, userID int64, tag string) ([]*model.Photo, error)
}
