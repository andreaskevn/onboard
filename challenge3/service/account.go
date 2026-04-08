package service

import (
	"challenge3/config"
	"challenge3/dto"
	"challenge3/models"
	"challenge3/repository"
	"challenge3/server"
	"context"

	"errors"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"net/http"
)

// type AccountService struct {
// 	repo *repository.AccountRepo
// }

//	func NewAccountService(repo *repository.AccountRepo) *AccountService {
//		return &AccountService{repo: repo}
//	}
type AccountService struct {
	repo     repository.IAccountRepo
	bankRepo repository.IBankRepo
}

func NewAccountService(repo repository.IAccountRepo, bankRepo repository.IBankRepo) *AccountService {
	return &AccountService{
		repo:     repo,
		bankRepo: bankRepo,
	}
}

func (t *AccountService) GetAllAccount() ([]models.Account, error) {
	data, err := t.repo.GetAll()
	if err != nil {
		config.Log.Error("Failed get all accounts",
			zap.Error(err),
		)

		return nil, err
	}
	config.Log.Info("Successfully get data all accounts") // zap.String(err),

	return data, nil
}

func (t *AccountService) GetAccountById(ctx context.Context, id string) (*models.Account, error) {
	data, err := t.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		config.Log.Error("Failed get data by acc id", zap.Error(err))
		return nil, err
	}
	config.Log.Info("Successfully get data by id",
		zap.String("account_id", data.ID.String()),
		zap.String("account_holder", data.AccountHolder),
		zap.String("account_number", data.AccountNumber),
		zap.Int("balance", data.Balance),
		zap.String("bank_id", data.BankID.String()),
	) // zap.String(err),
	return data, nil
}

func (t *AccountService) GetAccountByAccNumber(ctx context.Context, accNumber string) (*models.Account, error) {
	context, span := tracer.Start(ctx, "account.service.get-account-by-number")
	defer span.End()

	data, err := t.repo.GetByAccountNumber(context, accNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		config.Log.Error("Failed get data by acc number", zap.Error(err))
		return nil, err
	}

	config.Log.Info("Successfully get data by account number",
		zap.String("account_id", data.ID.String()),
		zap.String("account_holder", data.AccountHolder),
		zap.String("account_number", data.AccountNumber),
		zap.Int("balance", data.Balance),
		zap.String("bank_id", data.BankID.String()),
	) // zap.String(err),
	return data, nil
}

var tracer trace.Tracer = otel.Tracer("bank-service")

func (t *AccountService) CreateAcc(ctx context.Context, acc *models.Account) (*models.Account, error) {
	context, span := tracer.Start(ctx, "account.service.create-acc")
	defer span.End()
	_, err := t.bankRepo.GetById(context, acc.BankID.String())
	if err != nil {
		config.Log.Error("Failed create account bcs Bank ID doesnt exists",
			zap.Error(err),
		)

		server.AccountFailedCounter.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("reason", "Failed create account bcs Bank ID doesnt exists"),
			),
		)
		return nil, &dto.ErrorResponse{
			Message: "Bank ID doesnt exists",
			Code:    http.StatusNotFound,
		}
	}

	existing, err := t.GetAccountByAccNumber(context, acc.AccountNumber)
	if err == nil && existing != nil {
		// return nil, errors.New("Account with this number already exists")
		config.Log.Error("Failed create account bcs Account with this number already exists",
			zap.Error(err),
		)

		server.AccountFailedCounter.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("reason", "Failed create account bcs Account with this number already exists"),
			),
		)

		return nil, &dto.ErrorResponse{
			Message: "Account with this number already exists",
			Code:    http.StatusConflict,
		}
	}

	if acc.Balance <= 0 {
		config.Log.Error("Balance cant be lower than 0",
			zap.Error(err),
		)

		server.AccountFailedCounter.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("reason", "Balance cant be lower than 0"),
			),
		)

		return nil, &dto.ErrorResponse{
			Message: "Balance cant be lower than 0",
			Code:    http.StatusBadRequest,
		}
		// return nil, errors.New("Balance cant be lower than 0")
	}

	create, err := t.repo.CreateAcc(context, acc)
	if err != nil {
		config.Log.Error("Failed create account",
			zap.Error(err),
		)

		server.AccountFailedCounter.Add(
			ctx,
			1,
			metric.WithAttributes(
				attribute.String("reason", err.Error()),
			),
		)

		return nil, err
	}

	config.Log.Info("Successfully get data by account number",
		zap.String("account_id", create.ID.String()),
		zap.String("account_holder", create.AccountHolder),
		zap.String("account_number", create.AccountNumber),
		zap.Int("balance", create.Balance),
		zap.String("bank_id", create.BankID.String()),
	) // zap.String(err),

	server.AccountCreatedCounter.Add(ctx, 1)
	return create, nil
}

func (t *AccountService) UpdateAcc(ctx context.Context, id uuid.UUID, accountHolder string, balance int) (*models.Account, error) {
	existing, err := t.GetAccountById(ctx, id.String())
	if err != nil || existing == nil {
		return nil, &dto.ErrorResponse{
			Message: "Account not found",
			Code:    http.StatusNotFound,
		}
	}

	if existing.AccountHolder != accountHolder {
		return nil, &dto.ErrorResponse{
			Message: "Account Holder not Match with the Data",
			Code:    http.StatusBadRequest,
		}
	}

	if balance <= 0 {
		return nil, &dto.ErrorResponse{
			Message: "Balance cant be lower than 0",
			Code:    http.StatusBadRequest,
		}
	}

	return t.repo.UpdateAcc(id, accountHolder, balance)
}

func (t *AccountService) DeleteAcc(ctx context.Context, id string) error {
	existing, err := t.repo.GetById(ctx, id)
	if err != nil || existing == nil {
		return &dto.ErrorResponse{
			Message: "Account with this id doesnt exists",
			Code:    http.StatusNotFound,
		}
	}

	return t.repo.DeleteAcc(id)
}
