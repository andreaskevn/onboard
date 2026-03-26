package service

import (
	"challenge3/models"
	"challenge3/repository"

	"github.com/google/uuid"
	// "github.com/google/uuid"
	// "gorm.io/gorm"
)

type BankService struct {
	repo *repository.BankRepo
}

func NewBankService(repo *repository.BankRepo) *BankService {
	return &BankService{repo: repo}
}

func (t *BankService) GetAllBank() ([]models.Bank, error) {
	return t.repo.GetAll()
}

func (t *BankService) GetById(id string) (*models.Bank, error) {
	return t.repo.GetById(id)
}

func (t *BankService) GetByCode(code string) (*models.Bank, error) {
	return t.repo.GetByCode(code)
}

func (t *BankService) GetByName(name string) (*models.Bank, error) {
	return t.repo.GetByName(name)
}

func (t *BankService) CreateBank(bank *models.Bank) (*models.Bank, error) {
	return t.repo.CreateBank(bank)
}

func (t *BankService) UpdateBank(id uuid.UUID, bankName string, bankCode string) (*models.Bank, error) {
	return t.repo.UpdateBank(id, bankName, bankCode)
}

func (t *BankService) DeleteBank(id string) (error) {
	return t.repo.DeleteBank(id)
}
