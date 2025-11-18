package service

import (
	"context"
	"picstagsbot/internal/domain/model"
	"picstagsbot/internal/domain/repo"
	apperrors "picstagsbot/pkg/errors"
	"picstagsbot/pkg/logx"
	"picstagsbot/pkg/validator"
	"time"
)

type RegService struct {
	userRepo repo.UserRepo
}

func NewRegService(userRepo repo.UserRepo) *RegService {
	rs := &RegService{}

	rs.userRepo = userRepo

	return rs
}

func (svc *RegService) RegisterUser(ctx context.Context, telegramID int64, username string) (bool, error) {
	username = validator.SanitizeString(username)
	if err := validator.ValidateUsername(username); err != nil {
		logx.Warn("invalid username", "telegram_id", telegramID, "username", username, "error", err)
		return false, apperrors.ValidationError(err.Error())
	}

	existingUser, err := svc.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		logx.Error("failed to get user", "telegram_id", telegramID, "error", err)
		return false, apperrors.DatabaseError("failed to get user", err)
	}

	if existingUser != nil {
		logx.Info("user already registered", "telegram_id", telegramID, "username", username)
		return true, nil
	}

	botUser := &model.User{
		TelegramID: telegramID,
		Username:   username,
		CreatedAt:  time.Now(),
	}

	if err := svc.userRepo.Create(ctx, botUser); err != nil {
		logx.Error("failed to create user", "telegram_id", telegramID, "username", username, "error", err)
		return false, apperrors.DatabaseError("failed to create user", err)
	}

	logx.Info("user registered", "telegram_id", telegramID, "username", username, "user_id", botUser.ID)
	return false, nil
}
