package bootstrap

import (
	"database/sql"
	"net/http"

	userService "go_pg_http/internal/application/user/service"
	appConfig "go_pg_http/internal/infrastructure/config"
)

type App struct {
	Config      appConfig.Config
	DB          *sql.DB
	UserService *userService.UserService
	Router      http.Handler
}

func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
