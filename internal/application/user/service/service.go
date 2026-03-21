package service

import (
	"context"
	"go_pg_http/internal/application/user/dto"
	"go_pg_http/internal/application/user/ports"
	"go_pg_http/internal/application/user/usecase"
)

type UserService struct {
	createUserUseCase    *usecase.CreateUserUseCase
	getUserByNameUseCase *usecase.GetUserByNameUseCase
	listUsersUseCase     *usecase.ListUsersUseCase
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{
		createUserUseCase:    usecase.NewCreateUserUseCase(repo),
		getUserByNameUseCase: usecase.NewGetUserByNameUseCase(repo),
		listUsersUseCase:     usecase.NewListUsersUseCase(repo),
	}
}

func (s *UserService) Create(ctx context.Context, input dto.CreateUserInput) (dto.UserOutput, error) {
	return s.createUserUseCase.Execute(ctx, input)
}

func (s *UserService) GetByName(ctx context.Context, input dto.GetUserByNameInput) (dto.UserOutput, error) {
	return s.getUserByNameUseCase.Execute(ctx, input)
}

func (s *UserService) List(ctx context.Context, input dto.ListUsersInput) (dto.ListUsersOutput, error) {
	return s.listUsersUseCase.Execute(ctx, input)
}
