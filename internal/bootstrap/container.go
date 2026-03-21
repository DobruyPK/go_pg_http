package bootstrap

import (
	"context"
	"fmt"
)

func BuildApp(ctx context.Context) (*App, error) {
	cfg, err := ProvideConfig()
	if err != nil {
		return nil, fmt.Errorf("provide config: %w", err)
	}

	db, err := ProvideDB(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("provide db: %w", err)
	}

	userSvc := ProvideUserService(db)
	userHandler := ProvideUserHandler(userSvc)
	router := ProvideRouter(userHandler)

	return &App{
		Config:      cfg,
		DB:          db,
		UserService: userSvc,
		Router:      router,
	}, nil
}
