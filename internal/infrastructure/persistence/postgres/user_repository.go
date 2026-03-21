package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"

	applicationPorts "go_pg_http/internal/application/user/ports"
	domainUser "go_pg_http/internal/domain/user"
)

const (
	createUserQuery = `
		INSERT INTO users (name)
		VALUES ($1)
		RETURNING id, name
	`

	getUserByNameQuery = `
		SELECT id, name
		FROM users
		WHERE name = $1
	`

	listUsersQuery = `
		SELECT id, name
		FROM users
		ORDER BY id
	`
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

var _ applicationPorts.UserRepository = (*UserRepository)(nil)

func (r *UserRepository) Create(ctx context.Context, usr domainUser.User) (domainUser.User, error) {
	row := toUserRow(usr)

	var savedRow UserRow
	err := r.db.QueryRowContext(ctx, createUserQuery, row.Name).Scan(
		&savedRow.ID,
		&savedRow.Name,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 = unique_violation
			if pgErr.Code == "23505" {
				return domainUser.User{}, domainUser.ErrUserExists
			}
		}

		return domainUser.User{}, fmt.Errorf("create user: %w", err)
	}

	savedUser, err := toDomainUser(savedRow)
	if err != nil {
		return domainUser.User{}, fmt.Errorf("map saved user to domain: %w", err)
	}

	return savedUser, nil
}

func (r *UserRepository) GetByName(ctx context.Context, name domainUser.Name) (domainUser.User, error) {
	var row UserRow

	err := r.db.QueryRowContext(ctx, getUserByNameQuery, name.String()).Scan(
		&row.ID,
		&row.Name,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domainUser.User{}, domainUser.ErrUserNotFound
		}

		return domainUser.User{}, fmt.Errorf("get user by name: %w", err)
	}

	usr, err := toDomainUser(row)
	if err != nil {
		return domainUser.User{}, fmt.Errorf("map user row to domain: %w", err)
	}

	return usr, nil
}

func (r *UserRepository) List(ctx context.Context) ([]domainUser.User, int64, error) {
	rows, err := r.db.QueryContext(ctx, listUsersQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("list users query: %w", err)
	}
	defer rows.Close()

	users := make([]domainUser.User, 0)

	for rows.Next() {
		var row UserRow

		if err := rows.Scan(&row.ID, &row.Name); err != nil {
			return nil, 0, fmt.Errorf("scan user row: %w", err)
		}

		usr, err := toDomainUser(row)
		if err != nil {
			return nil, 0, fmt.Errorf("map listed user to domain: %w", err)
		}

		users = append(users, usr)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate user rows: %w", err)
	}

	total := int64(len(users))

	return users, total, nil
}
