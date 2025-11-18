package config

import (
	"fmt"
	"os"
	"picstagsbot/pkg/constants"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Environment string

const (
	EnvDevelopment Environment = "development"
	EnvProduction  Environment = "production"
	EnvStaging     Environment = "staging"
)

type Config struct {
	Env Environment    `yaml:"environment"`
	TG  TGBotConfig    `yaml:"telegram"`
	PG  PostgresConfig `yaml:"postgres"`
	App AppConfig      `yaml:"app"`
}

type TGBotConfig struct {
	Token         string        `yaml:"token"`
	PollerTimeout time.Duration `yaml:"poller_timeout"`
}

type PostgresConfig struct {
	URL             string
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	Name            string        `yaml:"name"`
	SSLMode         string        `yaml:"sslmode"`
	MaxConns        int           `yaml:"max_conns"`
	MinConns        int           `yaml:"min_conns"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `yaml:"max_conn_idle_time"`
	ConnectTimeout  time.Duration `yaml:"connect_timeout"`
	QueryTimeout    time.Duration `yaml:"query_timeout"`
}

type AppConfig struct {
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	RequestTimeout  time.Duration `yaml:"request_timeout"`
}

func New(file string) (*Config, error) {
	cfg := &Config{}

	if file == ".env" {
		_ = godotenv.Load(".env")
		return loadFromEnv(cfg)
	}

	if file == "config.yaml" {
		return loadFromYAML(cfg, file)
	}

	return nil, fmt.Errorf("unknown config file: %s", file)
}

func (c *Config) IsDevelopment() bool {
	return c.Env == EnvDevelopment
}

func (c *Config) IsProduction() bool {
	return c.Env == EnvProduction
}

func loadFromYAML(cfg *Config, path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Set defaults if not provided
	setDefaults(cfg)

	cfg.PG.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.PG.User,
		cfg.PG.Password,
		cfg.PG.Host,
		cfg.PG.Port,
		cfg.PG.Name,
		cfg.PG.SSLMode,
	)

	return cfg, nil
}

func loadFromEnv(cfg *Config) (*Config, error) {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = string(EnvDevelopment)
	}
	cfg.Env = Environment(env)

	botToken := os.Getenv("BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("BOT_TOKEN is empty")
	}
	cfg.TG.Token = botToken

	host := os.Getenv("DB_HOST")
	portStr := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || portStr == "" || user == "" || pass == "" || dbName == "" {
		return nil, fmt.Errorf("postgres env vars are incomplete")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	cfg.PG.Host = host
	cfg.PG.Port = port
	cfg.PG.User = user
	cfg.PG.Password = pass
	cfg.PG.Name = dbName
	cfg.PG.SSLMode = sslmode

	setDefaults(cfg)

	cfg.PG.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		user, pass, host, port, dbName, sslmode,
	)

	return cfg, nil
}

func setDefaults(cfg *Config) {
	if cfg.Env == "" {
		cfg.Env = EnvDevelopment
	}

	if cfg.TG.PollerTimeout == 0 {
		cfg.TG.PollerTimeout = constants.BotPollerTimeout
	}

	if cfg.PG.MaxConns == 0 {
		cfg.PG.MaxConns = constants.DBMaxConns
	}
	if cfg.PG.MinConns == 0 {
		cfg.PG.MinConns = constants.DBMinConns
	}
	if cfg.PG.MaxConnLifetime == 0 {
		cfg.PG.MaxConnLifetime = constants.DBMaxConnLifetime
	}
	if cfg.PG.MaxConnIdleTime == 0 {
		cfg.PG.MaxConnIdleTime = constants.DBMaxConnIdleTime
	}
	if cfg.PG.ConnectTimeout == 0 {
		cfg.PG.ConnectTimeout = constants.DBConnectTimeout
	}
	if cfg.PG.QueryTimeout == 0 {
		cfg.PG.QueryTimeout = constants.DBQueryTimeout
	}

	if cfg.App.ShutdownTimeout == 0 {
		cfg.App.ShutdownTimeout = constants.ShutdownTimeout
	}
	if cfg.App.RequestTimeout == 0 {
		cfg.App.RequestTimeout = constants.DefaultRequestTimeout
	}
}
