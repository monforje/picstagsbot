package handler

import (
	"context"
	"picstagsbot/internal/service"
	"picstagsbot/internal/tg/keyboard"
	"picstagsbot/internal/tg/message"
	"picstagsbot/pkg/logx"
	"time"

	tele "gopkg.in/telebot.v4"
)

type RegHandler struct {
	regService *service.RegService
}

func NewRegHandler(regService *service.RegService) *RegHandler {
	rh := &RegHandler{}

	rh.regService = regService

	return rh
}

func (h *RegHandler) HandleRegister(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	telegramID := c.Sender().ID
	username := c.Sender().Username

	logx.Info("register request", "telegram_id", telegramID, "username", username)

	isExisting, err := h.regService.RegisterUser(ctx, telegramID, username)
	if err != nil {
		logx.Error("register failed", "telegram_id", telegramID, "error", err)
		return message.SendWithEmoji(c, message.EmojiRegisterError, message.MsgRegisterError)
	}

	if isExisting {
		return message.SendWithEmoji(c, message.EmojiWelcomeExisting, message.MsgWelcomeExisting, keyboard.MainMenu)
	}

	return message.SendWithEmoji(c, message.EmojiWelcomeNew, message.MsgWelcomeNew, keyboard.MainMenu)
}
