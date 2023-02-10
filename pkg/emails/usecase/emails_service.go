package internal

import (
	"applicationDesignTest/pkg/emails"
	"context"
)

var _ emails.Usecase = (*EmailsUsecase)(nil)

type EmailsUsecase struct{}

func NewEmailsUsecase() *EmailsUsecase {
	return &EmailsUsecase{}
}

func (s *EmailsUsecase) FindEmail(_ context.Context, in emails.FindEmailInDTO) (out emails.FindEmailOutDTO, err error) {
	out.Email.Email = in.FindParams.Email
	out.Email.Status = emails.Confirmed
	return
}
