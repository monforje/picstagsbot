package repo

import (
	"context"
	"picstagsbot/internal/domain/model"
)

type UserRepo interface {
	Create(ctx context.Context, botuser *model.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
}
