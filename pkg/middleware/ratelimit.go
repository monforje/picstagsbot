package middleware

import (
	"picstagsbot/pkg/constants"
	"picstagsbot/pkg/logx"
	"sync"
	"time"

	tele "gopkg.in/telebot.v4"
)

type RateLimiter struct {
	mu       sync.RWMutex
	limits   map[int64]*userLimit
	requests int
	window   time.Duration
}

type userLimit struct {
	count     int
	resetTime time.Time
}

func NewRateLimiter(requests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limits:   make(map[int64]*userLimit),
		requests: requests,
		window:   window,
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) Middleware() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			userID := c.Sender().ID

			if !rl.Allow(userID) {
				logx.Warn("rate limit exceeded", "telegram_id", userID)
				return c.Send("Слишком много запросов. Пожалуйста, подождите немного.")
			}

			return next(c)
		}
	}
}

func (rl *RateLimiter) Allow(userID int64) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	limit, exists := rl.limits[userID]

	if !exists || now.After(limit.resetTime) {
		rl.limits[userID] = &userLimit{
			count:     1,
			resetTime: now.Add(rl.window),
		}
		return true
	}

	if limit.count >= rl.requests {
		return false
	}

	limit.count++
	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(constants.RateLimitCleanup)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for userID, limit := range rl.limits {
			if now.After(limit.resetTime) {
				delete(rl.limits, userID)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.limits = make(map[int64]*userLimit)
}
