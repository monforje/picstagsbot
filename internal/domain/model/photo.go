package model

import "time"

type Photo struct {
	ID          int64
	UserID      int64
	TelegramID  string
	FileSize    int64
	Width       int
	Height      int
	Description string
	Tags        []string
	CreatedAt   time.Time
}
