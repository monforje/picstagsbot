package search

import (
	"context"
	"fmt"
	"picstagsbot/internal/tg/keyboard"
	"picstagsbot/internal/tg/message"
	"picstagsbot/pkg/constants"
	"picstagsbot/pkg/logx"

	tele "gopkg.in/telebot.v4"
)

func (h *SearchHandler) HandleSearchStart(c tele.Context) error {
	userID := c.Sender().ID
	logx.Info("search started", "telegram_id", userID)
	h.setSession(userID)
	return message.SendWithEmoji(c, message.EmojiSearchPrompt, message.MsgSearchPrompt)
}

func (h *SearchHandler) HandleSearchQuery(c tele.Context, tag string) error {
	userID := c.Sender().ID

	if !h.getSession(userID) {
		return nil
	}

	logx.Info("search query", "telegram_id", userID, "tag", tag)

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	photos, err := h.searchService.SearchPhotosByTag(ctx, userID, tag)

	h.clearSession(userID)

	if err != nil {
		logx.Error("search failed", "telegram_id", userID, "tag", tag, "error", err)
		return message.SendWithEmoji(c, message.EmojiSearchError, message.MsgSearchError, keyboard.MainMenu)
	}

	if len(photos) == 0 {
		logx.Info("search no results", "telegram_id", userID, "tag", tag)
		return message.SendWithEmoji(c, message.EmojiSearchNoResults, message.MsgSearchNoResults, keyboard.MainMenu)
	}

	logx.Info("search completed", "telegram_id", userID, "tag", tag, "results_count", len(photos))

	resultMsg := fmt.Sprintf(message.MsgSearchResults, len(photos))
	if err := message.SendWithEmoji(c, message.EmojiSearchResults, resultMsg); err != nil {
		return err
	}

	h.sendPhotosAsAlbums(c, userID, photos)

	return message.SendWithEmoji(c, message.EmojiSearchCompleted, message.MsgSearchCompleted, keyboard.MainMenu)
}

func (h *SearchHandler) IsSearchingSession(userID int64) bool {
	return h.getSession(userID)
}
