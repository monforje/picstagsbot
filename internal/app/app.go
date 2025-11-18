package app

import (
	"context"
	"picstagsbot/config"
	"picstagsbot/internal/postgres"
	"picstagsbot/internal/postgres/repoimpl"
	"picstagsbot/internal/service"
	"picstagsbot/internal/tg/bot"
	"picstagsbot/internal/tg/handler"
	"picstagsbot/internal/tg/router"
	"picstagsbot/pkg/logx"
	"picstagsbot/pkg/middleware"
	"sync"
	"time"
)

type App struct {
	bot         *bot.Bot
	pg          *postgres.Postgres
	router      *router.Router
	rateLimiter *middleware.RateLimiter
	cfg         *config.Config
	wg          sync.WaitGroup
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	cfg, err := config.New(".env")
	if err != nil {
		return nil, err
	}
	a.cfg = cfg

	pgConfig := postgres.Config{
		URL:             cfg.PG.URL,
		MaxConns:        int32(cfg.PG.MaxConns),
		MinConns:        int32(cfg.PG.MinConns),
		MaxConnLifetime: cfg.PG.MaxConnLifetime,
		MaxConnIdleTime: cfg.PG.MaxConnIdleTime,
		ConnectTimeout:  cfg.PG.ConnectTimeout,
	}

	pg, err := postgres.New(pgConfig)
	if err != nil {
		return nil, err
	}
	a.pg = pg

	repo := repoimpl.New(pg.Pool)
	svc := service.New(repo)
	h := handler.New(svc)

	b, err := bot.New(cfg.TG.Token, cfg.TG.PollerTimeout)
	if err != nil {
		pg.Stop()
		return nil, err
	}
	a.bot = b

	rateLimiter := middleware.NewRateLimiter(20, 1*time.Minute)
	a.rateLimiter = rateLimiter

	r := router.New(b.Bot(), h, rateLimiter)
	a.router = r

	logx.Info("app initialized", "environment", cfg.Env)

	return a, nil
}

func (a *App) Run(ctx context.Context) {
	logx.Info("app starting...")

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		a.bot.Start()
	}()

	<-ctx.Done()
	a.Stop()
}

func (a *App) Stop() {
	logx.Info("app stopping")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.App.ShutdownTimeout)
	defer cancel()

	done := make(chan struct{})

	go func() {
		if a.bot != nil {
			a.bot.Stop()
		}

		a.wg.Wait()

		if a.pg != nil {
			a.pg.Stop()
		}

		close(done)
	}()

	select {
	case <-done:
		logx.Info("app stopped gracefully")
	case <-shutdownCtx.Done():
		logx.Warn("app shutdown timeout exceeded, forcing stop")
	}
}
