package bootstrap

import (
	"context"
	"database/sql"
	"net/http"

	userService "go_pg_http/internal/application/user/service"
	appConfig "go_pg_http/internal/infrastructure/config"
	httpInfra "go_pg_http/internal/infrastructure/http"
	"go_pg_http/internal/infrastructure/http/handlers"
	"go_pg_http/internal/infrastructure/persistence/postgres"
)

func ProvideConfig() (appConfig.Config, error) {
	return appConfig.Load()
}

func ProvideDB(ctx context.Context, cfg appConfig.Config) (*sql.DB, error) {
	return postgres.NewDB(ctx, postgres.DBConfig{
		Host:            cfg.Postgres.Host,
		Port:            cfg.Postgres.Port,
		User:            cfg.Postgres.User,
		Password:        cfg.Postgres.Password,
		Database:        cfg.Postgres.Database,
		SSLMode:         cfg.Postgres.SSLMode,
		MaxOpenConns:    cfg.Postgres.MaxOpenConns,
		MaxIdleConns:    cfg.Postgres.MaxIdleConns,
		ConnMaxLifetime: cfg.Postgres.ConnMaxLifetime,
	})
}

func ProvideUserService(db *sql.DB) *userService.UserService {
	repo := postgres.NewUserRepository(db)
	return userService.NewUserService(repo)
}

func ProvideUserHandler(userSvc *userService.UserService) *handlers.UserHandler {
	return handlers.NewUserHandler(userSvc)
}

func ProvideRouter(userHandler *handlers.UserHandler) http.Handler {
	return httpInfra.NewRouter(userHandler)
}
