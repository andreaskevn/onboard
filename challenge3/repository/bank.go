package repository

import (
	"challenge3/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	// "fmt"
)

type BankRepo struct {
	db *gorm.DB
}

func NewBankRepo(db *gorm.DB) *BankRepo {
	return &BankRepo{db: db}
}

func (t *BankRepo) GetAll() ([]models.Bank, error) {
	var banks []models.Bank

	err := t.db.Find(&banks).Error
	if err != nil {
		return nil, err
	}

	return banks, nil
}

func (t *BankRepo) GetById(id string) (*models.Bank, error) {
	var bank models.Bank

	err := t.db.First(&bank, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (t *BankRepo) GetByCode(code string) (*models.Bank, error) {
	var bank models.Bank

	err := t.db.First(&bank, "code = ?", code).Error
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (t *BankRepo) GetByName(name string) (*models.Bank, error) {
	var bank models.Bank

	err := t.db.First(&bank, "name = ?", name).Error
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (t *BankRepo) CreateBank(bank *models.Bank) (*models.Bank, error) {
	bank.ID = uuid.New()
	err := t.db.Create(bank).Error

	if err != nil {
		return nil, err
	}

	if err := t.db.Preload("Account").First(bank, "id = ?", bank.ID).Error; err != nil {
		return nil, err
	}

	return bank, nil
}

func (t *BankRepo) UpdateBank(id uuid.UUID, bankName string, bankCode string) (*models.Bank, error) {
	var bank models.Bank

	err := t.db.First(&bank, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	// if acc.

	bank.Code = bankCode
	bank.Name = bankName

	err = t.db.Save(&bank).Error
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

func (t *BankRepo) DeleteBank(id string) error {
	var bank models.Bank

	err := t.db.Delete(&bank, "id = ?", id).Error
	if err != nil {
		return err
	}

	return nil
}
