package usecase

import (
	"context"
	"go_pg_http/internal/application/user/dto"
	"go_pg_http/internal/application/user/ports"
)

type ListUsersUseCase struct {
	repo ports.UserRepository
}

func NewListUsersUseCase(repo ports.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{repo: repo}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, _ dto.ListUsersInput) (dto.ListUsersOutput, error) {
	users, total, err := uc.repo.List(ctx)
	if err != nil {
		return dto.ListUsersOutput{}, err
	}
	outputUsers := make([]dto.UserOutput, 0, len(users))
	for _, usr := range users {
		outputUsers = append(outputUsers, dto.UserOutput{
			ID:   usr.ID,
			Name: usr.Name.String(),
		})
	}
	return dto.ListUsersOutput{
		Users: outputUsers,
		Total: total,
	}, nil
}
