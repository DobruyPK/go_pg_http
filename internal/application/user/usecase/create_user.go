package usecase

import (
	"context"
	"go_pg_http/internal/application/user/dto"
	"go_pg_http/internal/application/user/ports"
	domainUser "go_pg_http/internal/domain/user"
)

// это не просто передача типаа в interface в usecase  что бы его могло что то реализовать посути это тупо передача обьекта что бы не тянуть ссылки
// потом туда будет запихнуто что то что реализиует этот иинтерфейс сейчас тут просто сам инрфейс без реализации
type CreateUserUseCase struct {
	repo ports.UserRepository
}

func NewCreateUserUseCase(repo ports.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo: repo,
	}
}

// метод типа CreateUserUseCase посутти точка вхта в usercase которая нам вернет отбьект dto используя реализацию интерфейса
func (uc *CreateUserUseCase) Execute(ctx context.Context, input dto.CreateUserInput) (dto.UserOutput, error) {
	usr, err := domainUser.New(0, input.Name)
	if err != nil {
		return dto.UserOutput{}, err
	}
	createUser, err := uc.repo.Create(ctx, usr)
	if err != nil {
		return dto.UserOutput{}, err
	}
	return dto.UserOutput{
		ID:   createUser.ID,
		Name: createUser.Name.String(),
	}, nil
}
