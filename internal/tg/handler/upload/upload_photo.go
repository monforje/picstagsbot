package upload

import (
	"context"
	"fmt"
	"picstagsbot/internal/tg/keyboard"
	"picstagsbot/internal/tg/message"
	"picstagsbot/pkg/constants"
	"picstagsbot/pkg/logx"
	"picstagsbot/pkg/validator"
	"time"

	tele "gopkg.in/telebot.v4"
)

func (h *UploadHandler) HandlePhoto(c tele.Context) error {
	userID := c.Sender().ID

	photo := c.Message().Photo
	if photo == nil {
		return message.SendWithEmoji(c, message.EmojiOnlyPhotoAllowed, message.MsgOnlyPhotoAllowed)
	}

	if err := validator.ValidateFileSize(int64(photo.FileSize)); err != nil {
		logx.Warn("photo file size exceeds limit", "telegram_id", userID, "file_size", photo.FileSize)
		return message.SendWithEmoji(c, message.EmojiPhotoSaveError, fmt.Sprintf("❌ Файл слишком большой. Максимальный размер: %d MB", constants.MaxFileSize/(1024*1024)))
	}

	session := h.getSession(userID)
	if session != nil {
		for _, p := range session.Photos {
			if p.FileID == photo.FileID {
				return message.SendWithEmoji(c, message.EmojiPhotoAlreadyExists, message.MsgPhotoAlreadyExists)
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	exists, err := h.uploadService.CheckPhotoExists(ctx, photo.FileID)
	if err != nil {
		h.clearSession(userID)
		logx.Error("photo check failed", "telegram_id", userID, "file_id", photo.FileID, "error", err)
		return message.SendWithEmoji(c, message.EmojiPhotoSaveError, message.MsgPhotoSaveError, keyboard.MainMenu)
	}

	if exists {
		return message.SendWithEmoji(c, message.EmojiPhotoAlreadyExists, message.MsgPhotoAlreadyExists)
	}

	newPhoto := UploadedPhoto{
		FileID:   photo.FileID,
		FileSize: int64(photo.FileSize),
		Width:    photo.Width,
		Height:   photo.Height,
	}

	if !h.addPhotoToSession(userID, newPhoto) {
		session := h.getSession(userID)
		if session != nil && len(session.Photos) >= constants.MaxPhotosPerSession {
			return message.SendWithEmoji(c, message.EmojiPhotoLimitReached, fmt.Sprintf(message.MsgPhotoLimitReached, constants.MaxPhotosPerSession))
		}
		return nil
	}

	logx.Info("photo added to session", "telegram_id", userID, "file_id", photo.FileID)

	mediaGroupID := c.Message().AlbumID

	if mediaGroupID != "" {
		h.setMediaGroupID(userID, mediaGroupID)

		time.AfterFunc(500*time.Millisecond, func() {
			h.sendPendingResponse(c, userID, mediaGroupID)
		})
		return nil
	}

	return message.SendWithEmoji(c, message.EmojiPhotoAdded, message.MsgPhotoAdded, keyboard.FinishUploadMenu)
}

func (h *UploadHandler) setMediaGroupID(userID int64, mediaGroupID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if session, ok := h.sessions[userID]; ok {
		session.LastMediaGroup = mediaGroupID
		session.PendingResponse = true
		session.LastActivity = time.Now()
	}
}

func (h *UploadHandler) sendPendingResponse(c tele.Context, userID int64, mediaGroupID string) {
	h.mu.Lock()
	session, ok := h.sessions[userID]
	if !ok || !session.PendingResponse || session.LastMediaGroup != mediaGroupID {
		h.mu.Unlock()
		return
	}
	session.PendingResponse = false
	h.mu.Unlock()

	message.SendWithEmoji(c, message.EmojiPhotoAdded, message.MsgPhotoAdded, keyboard.FinishUploadMenu)
}
