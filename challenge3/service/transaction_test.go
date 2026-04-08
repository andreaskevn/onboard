package service

import (
	"challenge3/models"
	"challenge3/repository/mocks"
	// "challenge3/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TransactionMocker struct {
	repo *mocks.ITransactionRepo
	accRepo *mocks.IAccountRepo
}

func TestTransactionService_GetHistory(t *testing.T) {

	id := uuid.New()

	testCases := []struct {
		desc      string
		mockSetup func(m *TransactionMocker)
		wantErr   bool
		expected  []models.Transaction
	}{
		{
			desc:    "SUCCESS: Get History",
			wantErr: false,
			mockSetup: func(m *TransactionMocker) {
				m.repo.On("GetHistory", id).
					Return([]models.Transaction{
						{Amount: 1000},
					}, nil)
			},
			expected: []models.Transaction{
				{Amount: 1000},
			},
		},
		{
			desc:    "ERROR: Failed Get History",
			wantErr: true,
			mockSetup: func(m *TransactionMocker) {
				m.repo.On("GetHistory", id).
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			m := &TransactionMocker{
				repo: mocks.NewITransactionRepo(t),
			}

			tC.mockSetup(m)

			svc := NewTransactionService(m.repo, m.accRepo)

			result, err := svc.GetHistory(t.Context(), id)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tC.expected, result)
			}
		})
	}
}

func TestTransactionService_Transfer(t *testing.T) {

	accFrom := uuid.New()
	accTo := uuid.New()

	testCases := []struct {
		desc      string
		mockSetup func(m *TransactionMocker)
		wantErr   bool
		expected  *models.Transaction
	}{
		{
			desc:    "SUCCESS: Transfer",
			wantErr: false,
			mockSetup: func(m *TransactionMocker) {
				m.repo.On("Transfer", accFrom, accTo, 1000, 500).
					Return(&models.Transaction{
						Amount: 1000,
					}, nil)
			},
			expected: &models.Transaction{
				Amount: 1000,
			},
		},
		{
			desc:    "ERROR: Transfer Failed",
			wantErr: true,
			mockSetup: func(m *TransactionMocker) {
				m.repo.On("Transfer", accFrom, accTo, 1000, 500).
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			m := &TransactionMocker{
				repo: mocks.NewITransactionRepo(t),
			}

			tC.mockSetup(m)

			svc := NewTransactionService(m.repo, m.accRepo)

			result, err := svc.Transfer(t.Context(), accFrom, accTo, 1000, 500)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tC.expected, result)
			}
		})
	}
}
