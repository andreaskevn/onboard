package repository

import (
	"challenge3/models"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	// "fmt"
)

type AccountRepo struct {
	db *gorm.DB
}

type IAccountRepo interface {
	GetAll() ([]models.Account, error)
	GetById(ctx context.Context, id string) (*models.Account, error)
	GetByAccountNumber(ctx context.Context, accNumber string) (*models.Account, error)
	CreateAcc(ctx context.Context, acc *models.Account) (*models.Account, error)
	UpdateAcc(id uuid.UUID, accountHolder string, balance int) (*models.Account, error)
	DeleteAcc(id string) error
}

func NewAccountRepo(db *gorm.DB) *AccountRepo {
	return &AccountRepo{db: db}
}

var tracer trace.Tracer = otel.Tracer("bank-service")

func (t *AccountRepo) GetAll() ([]models.Account, error) {
	var accounts []models.Account

	err := t.db.Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (t *AccountRepo) GetById(ctx context.Context, id string) (*models.Account, error) {
	var account models.Account

	_, span := tracer.Start(ctx, "account.repository.get-by-id")
	defer span.End()

	err := t.db.First(&account, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	if err := t.db.Preload("Bank").First(&account, "id = ?", account.ID).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (t *AccountRepo) GetByAccountNumber(ctx context.Context, accNumber string) (*models.Account, error) {
	var account models.Account

	_, span := tracer.Start(ctx, "account.repository.get-by-account-number")
	defer span.End()

	err := t.db.First(&account, "account_number = ?", accNumber).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (t *AccountRepo) CreateAcc(ctx context.Context, acc *models.Account) (*models.Account, error) {
	_, span := tracer.Start(ctx, "account.repository.create-account")
	defer span.End()

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
