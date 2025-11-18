package upload

import (
	"context"
	"picstagsbot/pkg/constants"
)

func (h *UploadHandler) savePhotosWithoutDescription(userID int64, photos []UploadedPhoto) int {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	savedCount := 0
	for _, p := range photos {
		isExisting, err := h.uploadService.UploadPhoto(
			ctx,
			userID,
			p.FileID,
			p.FileSize,
			p.Width,
			p.Height,
		)

		if err != nil {
			continue
		}

		if !isExisting {
			savedCount++
		}
	}
	return savedCount
}

func (h *UploadHandler) savePhotosWithDescription(userID int64, photos []UploadedPhoto, description string) int {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	savedCount := 0
	for _, p := range photos {
		_, err := h.uploadService.SavePhotoWithDescription(
			ctx,
			userID,
			p.FileID,
			p.FileSize,
			p.Width,
			p.Height,
			description,
		)

		if err != nil {
			continue
		}

		savedCount++
	}
	return savedCount
}
