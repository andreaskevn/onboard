package repository

import (
	"challenge3/models"
	// "challenge3/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "fmt"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(db *gorm.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (t *TransactionRepo) Transfer(accFrom uuid.UUID, accTo uuid.UUID, amount int) (*models.Transaction, error) {
	var tf models.Transaction

	tf.ID = uuid.New()
	tf.AccountFrom = accFrom
	tf.AccountTo = accTo
	tf.Amount = amount

	err := t.db.Create(&tf).Error
	if err != nil {
		return nil, err
	}

	return &tf, nil
}

func (t *TransactionRepo) GetHistory(id uuid.UUID) ([]models.Transaction, error) {
	var histories []models.Transaction

	err := t.db.
		Where("account_from = ? OR account_to = ?", id, id).
		Find(&histories).Error
	if err != nil {
		return nil, err
	}

	return histories, nil
}
