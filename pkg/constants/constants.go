package constants

import "time"

const (
	DBMaxConns           = 25
	DBMinConns           = 5
	DBMaxConnLifetime    = 1 * time.Hour
	DBMaxConnIdleTime    = 30 * time.Minute
	DBConnectTimeout     = 10 * time.Second
	DBQueryTimeout       = 30 * time.Second
	DBTransactionTimeout = 10 * time.Second
)

const (
	MaxFileSize         = 20 * 1024 * 1024
	MaxPhotosPerSession = 30
	MaxDescriptionLen   = 1000
	MaxTagLen           = 100
	MaxTagsPerPhoto     = 50
)

const (
	RateLimitRequests = 20
	RateLimitWindow   = 1 * time.Minute
	RateLimitCleanup  = 5 * time.Minute
)

const (
	SessionTimeout         = 30 * time.Minute
	SessionCleanupInterval = 10 * time.Minute
)

const (
	ShutdownTimeout       = 30 * time.Second
	BotPollerTimeout      = 10 * time.Second
	DefaultRequestTimeout = 15 * time.Second
)

const (
	TagPattern        = `^[a-zA-Z0-9а-яА-ЯёЁ_-]+$`
	MinTagLength      = 1
	MinUsernameLength = 1
	MaxUsernameLength = 255
)
