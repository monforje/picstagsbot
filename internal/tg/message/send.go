package message

import tele "gopkg.in/telebot.v4"

func SendWithEmoji(c tele.Context, emoji, text string, options ...interface{}) error {
	if err := c.Send(emoji); err != nil {
		return err
	}

	return c.Send(text, options...)
}
