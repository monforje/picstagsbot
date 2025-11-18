package router

import (
	"picstagsbot/internal/tg/handler"
	"picstagsbot/internal/tg/keyboard"
	"picstagsbot/internal/tg/message"
	"picstagsbot/pkg/logx"
	"picstagsbot/pkg/middleware"

	tele "gopkg.in/telebot.v4"
)

type Router struct {
	handler *handler.Handler
}

func New(b *tele.Bot, h *handler.Handler, rateLimiter *middleware.RateLimiter) *Router {
	r := &Router{handler: h}

	b.Use(rateLimiter.Middleware())

	b.Handle("/start", h.Reg.HandleRegister)
	b.Handle("/help", h.Help.HandleHelp)
	b.Handle("/info", h.Info.HandleInfo)

	b.Handle(&keyboard.BtnUploadPhoto, h.Upload.HandleUploadStart)
	b.Handle(&keyboard.BtnFinishUpload, h.Upload.HandleFinishUpload)
	b.Handle(&keyboard.BtnAddDescription, h.Upload.HandleAddDescription)
	b.Handle(&keyboard.BtnSkipDescription, h.Upload.HandleSkipDescription)

	b.Handle(&keyboard.BtnSearchPhoto, h.Search.HandleSearchStart)

	b.Handle(tele.OnText, r.handleText)
	b.Handle(tele.OnPhoto, r.handlePhoto)

	logx.Info("router initialized with rate limiting")

	return r
}

func (r *Router) handleText(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return nil
	}

	userID := sender.ID
	text := c.Text()

	if len(text) > 0 && text[0] == '/' {
		return nil
	}

	if r.handler.Upload.IsUploadingSession(userID) {
		return r.handler.Upload.HandleText(c)
	}

	if r.handler.Search.IsSearchingSession(userID) {
		return r.handler.Search.HandleSearchQuery(c, text)
	}

	return message.SendWithEmoji(c, message.EmojiUseButtons, message.MsgUseButtons, keyboard.MainMenu)
}

func (r *Router) handlePhoto(c tele.Context) error {
	sender := c.Sender()
	if sender == nil {
		return nil
	}

	userID := sender.ID

	if r.handler.Upload.IsUploadingSession(userID) {
		return r.handler.Upload.HandlePhoto(c)
	}

	return message.SendWithEmoji(c, message.EmojiUseButtons, message.MsgUseButtons, keyboard.MainMenu)
}
