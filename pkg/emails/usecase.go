package emails

import (
	"context"
	"errors"
)

var ErrEmailNotExists = errors.New("email not exists")

type Usecase interface {
	// FindEmail должен находить информацию по email.
	//
	// Если его нет вернуть ErrEmailNotExists.
	FindEmail(ctx context.Context, in FindEmailInDTO) (out FindEmailOutDTO, err error)
}

type FindEmailInDTO struct {
	FindParams FindParamsInDTO
}

type FindParamsInDTO struct {
	Email string
}

type FindEmailOutDTO struct {
	Email EmailDTO
}
