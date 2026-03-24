package service

import (
	"challenge3/models"
	"challenge3/repository"

	"github.com/google/uuid"
	// "gorm.io/gorm"
)

type AccountService struct {
	repo *repository.AccountRepo
}

func NewAccountService(repo *repository.AccountRepo) *AccountService {
	return &AccountService{repo: repo}
}

func (t *AccountService) GetAllAccount() ([]models.Account, error) {
	return t.repo.GetAll()
} 

func (t *AccountService) GetAccountById(id string) (*models.Account, error) {
	return t.repo.GetById(id)
} 

func (t *AccountService) GetAccountByAccNumber(accNumber string) (*models.Account, error) {
	return t.repo.GetByAccountNumber(accNumber)
} 

func (t *AccountService) CreateAcc(acc *models.Account) (*models.Account, error) {
	return t.repo.CreateAcc(acc)
} 

func (t *AccountService) UpdateAcc(id uuid.UUID, accountHolder string, balance int) (*models.Account, error) {
	return t.repo.UpdateAcc(id, accountHolder, balance)
} 

func (t *AccountService) DeleteAcc(id string) (error) {
	return t.repo.DeleteAcc(id)
} 

