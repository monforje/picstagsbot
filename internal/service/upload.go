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

type UploadService struct {
	photoRepo repo.PhotoRepo
	userRepo  repo.UserRepo
}

func NewUploadService(photoRepo repo.PhotoRepo, userRepo repo.UserRepo) *UploadService {
	us := &UploadService{}

	us.photoRepo = photoRepo
	us.userRepo = userRepo

	return us
}

func (svc *UploadService) CheckPhotoExists(ctx context.Context, fileID string) (bool, error) {
	photo, err := svc.photoRepo.GetByFileID(ctx, fileID)
	if err != nil {
		logx.Error("failed to check photo existence", "file_id", fileID, "error", err)
		return false, apperrors.DatabaseError("failed to check photo existence", err)
	}
	return photo != nil, nil
}

func (svc *UploadService) UploadPhoto(ctx context.Context, telegramID int64, fileID string, fileSize int64, width, height int) (bool, error) {
	if err := validator.ValidateFileSize(fileSize); err != nil {
		logx.Warn("invalid file size", "telegram_id", telegramID, "size", fileSize, "error", err)
		return false, apperrors.ValidationError(err.Error())
	}

	existingPhoto, err := svc.photoRepo.GetByFileID(ctx, fileID)
	if err != nil {
		logx.Error("failed to check existing photo", "telegram_id", telegramID, "file_id", fileID, "error", err)
		return false, apperrors.DatabaseError("failed to check existing photo", err)
	}
	if existingPhoto != nil {
		logx.Info("photo already exists", "telegram_id", telegramID, "file_id", fileID)
		return true, nil
	}

	user, err := svc.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		logx.Error("failed to get user for photo upload", "telegram_id", telegramID, "error", err)
		return false, apperrors.DatabaseError("failed to get user", err)
	}
	if user == nil {
		logx.Warn("user not found for photo upload", "telegram_id", telegramID)
		return false, apperrors.NotFoundError("user not found")
	}

	photo := &model.Photo{
		UserID:     user.ID,
		TelegramID: fileID,
		FileSize:   fileSize,
		Width:      width,
		Height:     height,
		CreatedAt:  time.Now(),
	}

	if err := svc.photoRepo.Create(ctx, photo); err != nil {
		logx.Error("failed to create photo", "telegram_id", telegramID, "user_id", user.ID, "file_id", fileID, "error", err)
		return false, apperrors.DatabaseError("failed to create photo", err)
	}

	logx.Info("photo uploaded", "telegram_id", telegramID, "user_id", user.ID, "photo_id", photo.ID)
	return false, nil
}

func (svc *UploadService) SavePhotoWithDescription(ctx context.Context, telegramID int64, fileID string, fileSize int64, width, height int, description string) (*model.Photo, error) {
	if err := validator.ValidateFileSize(fileSize); err != nil {
		logx.Warn("invalid file size", "telegram_id", telegramID, "size", fileSize, "error", err)
		return nil, apperrors.ValidationError(err.Error())
	}

	tags, err := validator.ValidateAndParseTags(description)
	if err != nil {
		logx.Warn("invalid description or tags", "telegram_id", telegramID, "error", err)
		return nil, apperrors.ValidationError(err.Error())
	}

	description = validator.SanitizeString(description)

	existingPhoto, err := svc.photoRepo.GetByFileID(ctx, fileID)
	if err != nil {
		logx.Error("failed to check existing photo with description", "telegram_id", telegramID, "file_id", fileID, "error", err)
		return nil, apperrors.DatabaseError("failed to check existing photo", err)
	}
	if existingPhoto != nil {
		logx.Info("photo with description already exists", "telegram_id", telegramID, "file_id", fileID)
		return existingPhoto, nil
	}

	user, err := svc.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		logx.Error("failed to get user for photo with description", "telegram_id", telegramID, "error", err)
		return nil, apperrors.DatabaseError("failed to get user", err)
	}
	if user == nil {
		logx.Warn("user not found for photo with description", "telegram_id", telegramID)
		return nil, apperrors.NotFoundError("user not found")
	}

	photo := &model.Photo{
		UserID:      user.ID,
		TelegramID:  fileID,
		FileSize:    fileSize,
		Width:       width,
		Height:      height,
		Description: description,
		Tags:        tags,
		CreatedAt:   time.Now(),
	}

	if err := svc.photoRepo.Create(ctx, photo); err != nil {
		logx.Error("failed to create photo with description", "telegram_id", telegramID, "user_id", user.ID, "tags_count", len(tags), "error", err)
		return nil, apperrors.DatabaseError("failed to create photo", err)
	}

	logx.Info("photo saved with description", "telegram_id", telegramID, "user_id", user.ID, "photo_id", photo.ID, "tags_count", len(tags))
	return photo, nil
}

func (svc *UploadService) AddDescriptionToPhoto(ctx context.Context, photoID int64, description string) error {
	tags, err := validator.ValidateAndParseTags(description)
	if err != nil {
		logx.Warn("invalid description or tags", "photo_id", photoID, "error", err)
		return apperrors.ValidationError(err.Error())
	}

	description = validator.SanitizeString(description)

	err = svc.photoRepo.UpdateDescription(ctx, photoID, description, tags)
	if err != nil {
		logx.Error("failed to update photo description", "photo_id", photoID, "tags_count", len(tags), "error", err)
		return apperrors.DatabaseError("failed to update photo description", err)
	}

	logx.Info("photo description updated", "photo_id", photoID, "tags_count", len(tags))
	return nil
}
