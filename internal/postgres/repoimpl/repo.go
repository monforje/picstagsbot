package repoimpl

import (
	"picstagsbot/internal/domain/repo"
	"picstagsbot/pkg/logx"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(pool *pgxpool.Pool) *repo.Repo {
	r := &repo.Repo{}

	r.UserRepo = NewUserRepo(pool)
	r.PhotoRepo = NewPhotoRepo(pool)

	logx.Info("postgres repositories initialized")

	return r
}
