package bot

import (
	"picstagsbot/pkg/logx"
	"time"

	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	bot *tele.Bot
}

func New(token string, pollerTimeout time.Duration) (*Bot, error) {
	bot := &Bot{}

	tb, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: pollerTimeout},
	})
	if err != nil {
		return nil, err
	}
	bot.bot = tb

	logx.Info("telegram bot initialized", "poller_timeout", pollerTimeout)

	return bot, nil
}

func (b *Bot) Bot() *tele.Bot {
	return b.bot
}

func (b *Bot) Start() {
	logx.Info("telegram bot started")
	b.bot.Start()
}

func (b *Bot) Stop() {
	b.bot.Stop()
	logx.Info("telegram bot stopped")
}
