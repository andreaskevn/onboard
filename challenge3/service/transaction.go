package service

import (
	"challenge3/dto"
	"challenge3/models"
	"challenge3/repository"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

// type TransactionService struct {
// 	repo *repository.TransactionRepo
// }

// func NewTransactionService(repo *repository.TransactionRepo) *TransactionService {
// 	return &TransactionService{repo: repo}
// }

type TransactionService struct {
	repo    repository.ITransactionRepo
	accRepo repository.IAccountRepo
}

func NewTransactionService(repo repository.ITransactionRepo, accRepo repository.IAccountRepo) *TransactionService {
	return &TransactionService{
		repo:    repo,
		accRepo: accRepo,
	}
}

func (t *TransactionService) Transfer(ctx context.Context, accFrom uuid.UUID, accTo uuid.UUID, amount int, adminFee int) (*models.Transaction, error) {
	context, span := tracer.Start(ctx, "transacation.service.transfer")
	defer span.End()
	
	if accFrom == accTo {
		return nil, &dto.ErrorResponse{
			Message: "Account From and Destination cant be same",
			Code:    http.StatusBadRequest,
		}
	}

	accountFromExist, err := t.accRepo.GetById(context, accFrom.String())
	if err != nil || accountFromExist == nil {
		return nil, &dto.ErrorResponse{
			Message: "Account From doesnt exist",
			Code:    http.StatusNotFound,
		}
	}

	accountToExist, err := t.accRepo.GetById(context, accTo.String())
	if err != nil || accountToExist == nil {
		return nil, &dto.ErrorResponse{
			Message: "Account Destination doesnt exist",
			Code:    http.StatusNotFound,
		}
	}

	if amount <= 0 {
		return nil, &dto.ErrorResponse{
			Message: "Amount cant be lower than 0",
			Code:    http.StatusBadRequest,
		}
	}

	if accountFromExist.Balance < (amount + adminFee) {
		fmt.Print(accountFromExist)
		log.Print(accountFromExist)
		return nil, &dto.ErrorResponse{
			Message: "Your balance is insufficient",
			Code:    http.StatusNotFound,
		}
	}

	if accountFromExist.BankID != accountToExist.BankID {
		if adminFee <= 0 {
			return nil, &dto.ErrorResponse{
				Message: "Admin Fee must included",
				Code:    http.StatusNotFound,
			}
		}
	}

	reduceBalance := accountFromExist.Balance - (amount + adminFee)
	addBalance := accountToExist.Balance + amount

	_, err = t.accRepo.UpdateAcc(accountFromExist.ID, accountFromExist.AccountHolder, reduceBalance)
	if err != nil {
		return nil, err
	}
	_, err = t.accRepo.UpdateAcc(accountToExist.ID, accountToExist.AccountHolder, addBalance)
	if err != nil {
		return nil, err
	}
	// if err != nil {
	// 	return nil, &dto.ErrorResponse{
	// 		Message: "error",
	// 		Code:    http.StatusInternalServerError,
	// 	}
	// }

	return t.repo.Transfer(accFrom, accTo, amount, adminFee)
}

func (t *TransactionService) GetHistory(ctx context.Context, id uuid.UUID) ([]models.Transaction, error) {
	_, err := t.accRepo.GetById(ctx, id.String())
	if err != nil {
		return nil, &dto.ErrorResponse{
			Message: "Account doesnt exists",
			Code:    http.StatusNotFound,
		}
	}
	return t.repo.GetHistory(id)
}
