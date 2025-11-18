package upload

import (
	"picstagsbot/internal/tg/keyboard"
	"picstagsbot/internal/tg/message"
	"picstagsbot/pkg/constants"
	"picstagsbot/pkg/logx"

	tele "gopkg.in/telebot.v4"
)

func (h *UploadHandler) HandleUploadStart(c tele.Context) error {
	userID := c.Sender().ID

	logx.Info("upload started", "telegram_id", userID)

	h.setSession(userID, &UploadSession{
		State: StateAwaitingPhoto,
	})

	return message.SendWithEmoji(c, message.EmojiUploadPhotoPrompt, message.MsgUploadPhotoPrompt, keyboard.FinishUploadMenu)
}

func (h *UploadHandler) HandleAddDescription(c tele.Context) error {
	userID := c.Sender().ID

	session := h.getSession(userID)
	if session == nil || session.State != StateAwaitingDescription {
		return nil
	}

	return message.SendWithEmoji(c, message.EmojiEnterDescription, message.MsgEnterDescription)
}

func (h *UploadHandler) HandleFinishUpload(c tele.Context) error {
	userID := c.Sender().ID

	session := h.getSession(userID)
	if session == nil || len(session.Photos) == 0 {
		h.clearSession(userID)
		logx.Info("upload finished with no photos", "telegram_id", userID)
		return message.SendWithEmoji(c, message.EmojiNoPhotosToSave, message.MsgNoPhotosToSave, keyboard.MainMenu)
	}

	logx.Info("upload awaiting description", "telegram_id", userID, "photos_count", len(session.Photos))
	h.updateSessionState(userID, StateAwaitingDescription)

	return message.SendWithEmoji(c, message.EmojiPhotoReceived, message.MsgPhotoReceived, keyboard.DescriptionMenu)
}

func (h *UploadHandler) HandleSkipDescription(c tele.Context) error {
	userID := c.Sender().ID

	session := h.getSession(userID)
	if session == nil || session.State != StateAwaitingDescription {
		return nil
	}

	savedCount := h.savePhotosWithoutDescription(userID, session.Photos)

	h.clearSession(userID)

	if savedCount == 0 {
		logx.Error("upload failed - no photos saved", "telegram_id", userID)
		return message.SendWithEmoji(c, message.EmojiPhotoSaveError, message.MsgPhotoSaveError, keyboard.MainMenu)
	}

	logx.Info("upload completed without description", "telegram_id", userID, "saved_count", savedCount)
	return message.SendWithEmoji(c, message.EmojiPhotosSaved, message.MsgPhotosSaved, keyboard.MainMenu)
}

func (h *UploadHandler) HandleText(c tele.Context) error {
	userID := c.Sender().ID

	session := h.getSession(userID)
	if session == nil || session.State != StateAwaitingDescription {
		return nil
	}

	description := c.Text()

	if len(description) > 0 && description[0] == '/' {
		return nil
	}

	if len(description) > constants.MaxDescriptionLen {
		return message.SendWithEmoji(c, message.EmojiDescriptionTooLong, message.MsgDescriptionTooLong)
	}

	savedCount := h.savePhotosWithDescription(userID, session.Photos, description)

	h.clearSession(userID)

	if savedCount == 0 {
		logx.Error("upload failed - no photos saved with description", "telegram_id", userID)
		return message.SendWithEmoji(c, message.EmojiPhotoSaveError, message.MsgPhotoSaveError, keyboard.MainMenu)
	}

	logx.Info("upload completed with description", "telegram_id", userID, "saved_count", savedCount)
	return message.SendWithEmoji(c, message.EmojiPhotosSavedWithDesc, message.MsgPhotosSavedWithDesc, keyboard.MainMenu)
}

func (h *UploadHandler) IsUploadingSession(userID int64) bool {
	return h.getSession(userID) != nil
}
