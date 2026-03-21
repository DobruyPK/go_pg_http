package ports

import (
	"context"

	domainUser "go_pg_http/internal/domain/user"
)

type UserRepository interface {
	Create(ctx context.Context, usr domainUser.User) (domainUser.User, error)
	GetByName(ctx context.Context, name domainUser.Name) (domainUser.User, error)
	List(ctx context.Context) ([]domainUser.User, int64, error)
}
