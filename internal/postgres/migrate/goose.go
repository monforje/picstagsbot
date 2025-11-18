package migrate

import (
	"database/sql"
	"picstagsbot/pkg/logx"

	"github.com/pressly/goose/v3"
)

type Goose struct{}

func New() *Goose {
	g := &Goose{}

	goose.SetLogger(goose.NopLogger())

	logx.Info("goose migrator initialized")

	return g
}

func (g *Goose) Up(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		return err
	}

	logx.Info("migrations applied successfully")

	return nil
}
