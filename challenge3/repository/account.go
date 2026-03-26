package repository

import (
	"challenge3/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "fmt"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

func (t *AccountRepo) GetAll() ([]models.Account, error) {
	var accounts []models.Account

	err := t.db.Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (t *AccountRepo) GetById(id string) (*models.Account, error) {
	var account models.Account

	err := t.db.First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	if err := t.db.Preload("Bank").First(&account, "id = ?", account.ID).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (t *AccountRepo) GetByAccountNumber(accNumber string) (*models.Account, error) {
	var account models.Account

	err := t.db.First(&account, "account_number = ?", accNumber).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (t *AccountRepo) CreateAcc(acc *models.Account) (*models.Account, error) {
	acc.ID = uuid.New()
	err := t.db.Create(acc).Error

	if err != nil {
		return nil, err
	}

	if err := t.db.Preload("Bank").First(acc, "id = ?", acc.ID).Error; err != nil {
		return nil, err
	}

	return acc, nil
}

func (t *AccountRepo) UpdateAcc(id uuid.UUID, accountHolder string, balance int) (*models.Account, error) {
	var acc models.Account

	err := t.db.First(&acc, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// if acc.

	acc.AccountHolder = accountHolder
	acc.Balance = balance

	err = t.db.Save(&acc).Error
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

func (t *AccountRepo) DeleteAcc(id string) error {
	var acc models.Account

	err := t.db.Delete(&acc, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
