package service

import (
	"picstagsbot/internal/domain/repo"
	"picstagsbot/pkg/logx"
)

type Service struct {
	Reg    *RegService
	Upload *UploadService
	Search *SearchService
}

func New(repo *repo.Repo) *Service {
	s := &Service{}

	s.Reg = NewRegService(repo.UserRepo)
	s.Upload = NewUploadService(repo.PhotoRepo, repo.UserRepo)
	s.Search = NewSearchService(repo.PhotoRepo, repo.UserRepo)

	logx.Info("services initialized")

	return s
}
