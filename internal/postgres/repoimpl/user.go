package repoimpl

import (
	"context"
	"errors"
	"picstagsbot/internal/domain/model"
	"picstagsbot/pkg/logx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	ur := &UserRepo{}

	ur.pool = pool

	return ur
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (telegram_id, username, created_at) VALUES ($1, $2, $3) RETURNING id`
	err := r.pool.QueryRow(ctx, query, user.TelegramID, user.Username, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		logx.Error("db: failed to create user", "telegram_id", user.TelegramID, "error", err)
		return err
	}
	return nil
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	query := `SELECT id, telegram_id, username, created_at FROM users WHERE telegram_id = $1`

	row := r.pool.QueryRow(ctx, query, telegramID)

	var user model.User
	err := row.Scan(&user.ID, &user.TelegramID, &user.Username, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		logx.Error("db: failed to get user by telegram_id", "telegram_id", telegramID, "error", err)
		return nil, err
	}

	return &user, nil
}
