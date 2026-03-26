package service

import (
	"challenge3/models"
	"challenge3/repository"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type TransactionService struct {
	repo *repository.TransactionRepo
}

func NewTransactionService(repo *repository.TransactionRepo) *TransactionService {
	return &TransactionService{repo: repo}
}

func (t *TransactionService) Transfer(accFrom uuid.UUID, accTo uuid.UUID, amount int, adminFee int) (*models.Transaction, error) {
	return t.repo.Transfer(accFrom, accTo, amount, adminFee)
}

func (t *TransactionService) GetHistory(id uuid.UUID) ([]models.Transaction, error) {
	return t.repo.GetHistory(id)
}


