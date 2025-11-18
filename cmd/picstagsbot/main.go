package main

import (
	"context"
	"os"
	"os/signal"
	"picstagsbot/internal/app"
	"picstagsbot/pkg/logx"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	a, err := app.New(ctx)
	if err != nil {
		logx.Fatal("app init err: %s", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	a.Run(ctx)
}
