package usecase

import (
	"context"
	"go_pg_http/internal/application/user/dto"
	"go_pg_http/internal/application/user/ports"
	domainUser "go_pg_http/internal/domain/user"
)

type GetUserByNameUseCase struct {
	repo ports.UserRepository
}

func NewGetUserByNameUseCase(repo ports.UserRepository) *GetUserByNameUseCase {
	return &GetUserByNameUseCase{repo: repo}
}

func (uc *GetUserByNameUseCase) Execute(ctx context.Context, input dto.GetUserByNameInput) (dto.UserOutput, error) {
	name, err := domainUser.NewName(input.Name)
	if err != nil {
		return dto.UserOutput{}, err
	}
	usr, err := uc.repo.GetByName(ctx, name)
	if err != nil {
		return dto.UserOutput{}, err
	}
	return dto.UserOutput{
		ID:   usr.ID,
		Name: usr.Name.String(),
	}, nil
}
