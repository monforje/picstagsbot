package handler

import (
	"picstagsbot/internal/service"
	"picstagsbot/internal/tg/handler/search"
	"picstagsbot/internal/tg/handler/upload"
	"picstagsbot/pkg/logx"
)

type Handler struct {
	Reg    *RegHandler
	Help   *HelpHandler
	Info   *InfoHandler
	Upload *upload.UploadHandler
	Search *search.SearchHandler
}

func New(svc *service.Service) *Handler {
	h := &Handler{}

	h.Reg = NewRegHandler(svc.Reg)
	h.Help = NewHelpHandler()
	h.Info = NewInfoHandler()
	h.Upload = upload.NewUploadHandler(svc.Upload)
	h.Search = search.NewSearchHandler(svc.Search)

	logx.Info("handlers initialized")

	return h
}
