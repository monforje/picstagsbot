package handler

import (
	"picstagsbot/internal/tg/keyboard"
	"picstagsbot/internal/tg/message"

	tele "gopkg.in/telebot.v4"
)

type InfoHandler struct{}

func NewInfoHandler() *InfoHandler {
	info := &InfoHandler{}

	return info
}

func (h *InfoHandler) HandleInfo(c tele.Context) error {
	return c.Send(message.MsgInfo, keyboard.MainMenu)
}
