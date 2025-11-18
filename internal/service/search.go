package service

import (
	"context"
	"picstagsbot/internal/domain/model"
	"picstagsbot/internal/domain/repo"
	apperrors "picstagsbot/pkg/errors"
	"picstagsbot/pkg/logx"
	"picstagsbot/pkg/validator"
)

type SearchService struct {
	photoRepo repo.PhotoRepo
	userRepo  repo.UserRepo
}

func NewSearchService(photoRepo repo.PhotoRepo, userRepo repo.UserRepo) *SearchService {
	sh := &SearchService{}

	sh.photoRepo = photoRepo
	sh.userRepo = userRepo

	return sh
}

func (svc *SearchService) SearchPhotosByTag(ctx context.Context, telegramID int64, tag string) ([]*model.Photo, error) {
	tag = validator.SanitizeString(tag)
	if err := validator.ValidateTag(tag); err != nil {
		logx.Warn("invalid search tag", "telegram_id", telegramID, "tag", tag, "error", err)
		return nil, apperrors.ValidationError(err.Error())
	}

	user, err := svc.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		logx.Error("failed to get user for search", "telegram_id", telegramID, "tag", tag, "error", err)
		return nil, apperrors.DatabaseError("failed to get user", err)
	}

	if user == nil {
		logx.Warn("user not found for search", "telegram_id", telegramID, "tag", tag)
		return nil, apperrors.NotFoundError("user not found")
	}

	photos, err := svc.photoRepo.SearchByTag(ctx, user.ID, tag)
	if err != nil {
		logx.Error("failed to search photos by tag", "telegram_id", telegramID, "user_id", user.ID, "tag", tag, "error", err)
		return nil, apperrors.DatabaseError("failed to search photos", err)
	}

	logx.Info("photos searched by tag", "telegram_id", telegramID, "user_id", user.ID, "tag", tag, "results_count", len(photos))
	return photos, nil
}
