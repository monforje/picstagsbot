package main

import (
	"database/sql"
	"os"

	"picstagsbot/config"
	"picstagsbot/internal/postgres/migrate"
	"picstagsbot/pkg/logx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg, err := config.New("config.yaml")
	if err != nil {
		logx.Fatal("config init err: %s", err)
	}

	db, err := sql.Open("pgx", cfg.PG.URL)
	if err != nil {
		logx.Fatal("open db err: %s", err)
	}
	defer db.Close()

	if err := migrate.New().Up(db); err != nil {
		logx.Fatal("migrate up err: %s", err)
	}

	logx.Info("migrations applied successfully")

	_ = os.Stdout
}
