package handler

import (
	"picstagsbot/internal/tg/keyboard"
	"picstagsbot/internal/tg/message"

	tele "gopkg.in/telebot.v4"
)

type HelpHandler struct{}

func NewHelpHandler() *HelpHandler {
	hh := &HelpHandler{}

	return hh
}

func (h *HelpHandler) HandleHelp(c tele.Context) error {
	return c.Send(message.MsgHelp, keyboard.MainMenu)
}
