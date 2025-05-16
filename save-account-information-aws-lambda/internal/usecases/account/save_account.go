package usecases

import (
	"context"
	"fmt"
	"stori-card-challenge/save-account-information-aws-lambda/domain/account"
	infraAccount "stori-card-challenge/save-account-information-aws-lambda/internal/infrastructure/account"

	"github.com/pkg/errors"
)

type SaveAccountUsecase interface {
	Execute(ctx context.Context, account *account.Account) error
}

type saveAccountUsecase struct {
	accountRepository infraAccount.AccountDBRepository
}

func NewSaveAccountUsecase(accountRepository infraAccount.AccountDBRepository) *saveAccountUsecase {
	return &saveAccountUsecase{
		accountRepository: accountRepository,
	}
}

func (s *saveAccountUsecase) Execute(ctx context.Context, account *account.Account) error {
	err := validateModel(account)

	if err != nil {
		return errors.Wrap(err, "model is not defined correctly")
	}
	err = s.accountRepository.SaveUserAccount(ctx, account)

	if err != nil {
		return errors.Wrap(err, " cannot save user account")

	}
	return err
}

func validateModel(account *account.Account) error {
	if account.Id == "" {
		return fmt.Errorf("ID is required")
	}
	if account.DateCreated.IsZero() {
		return fmt.Errorf("DateCreated is required")
	}

	return nil
}
