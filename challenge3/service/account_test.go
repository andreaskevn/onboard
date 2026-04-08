package service

import (
	"challenge3/models"
	"challenge3/repository/mocks"
	// "challenge3/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Mocker struct {
	repo     *mocks.IAccountRepo
	bankRepo *mocks.IBankRepo
}

func TestAccountService_GetAccountById(t *testing.T) {

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
		expected  *models.Account
	}{
		{
			desc:    "SUCCESS: Get Account by ID",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetById",
					mock.Anything,
					mock.AnythingOfType("string"),
				).Return(&models.Account{
					AccountNumber: "123",
					AccountHolder: "John Doe",
					Balance:       100000,
				}, nil)
			},
			expected: &models.Account{
				AccountNumber: "123",
				AccountHolder: "John Doe",
				Balance:       100000,
			},
		},
		{
			desc:    "ERROR: Account not found",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On(
					"GetById",
					mock.Anything,
					mock.AnythingOfType("string"),
				).Return(nil, assert.AnError)
			},
			expected: nil,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			m := &Mocker{
				repo: mocks.NewIAccountRepo(t),
				// bankRepo: mocks.NewIBankRepo(t),
			}

			tC.mockSetup(m)

			service := NewAccountService(m.repo, m.bankRepo)

			result, err := service.GetAccountById(t.Context(), "123")

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tC.expected, result)
			}

			m.repo.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAllAccount(t *testing.T) {
	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
		expected  []models.Account
	}{
		{
			desc:    "SUCCESS: Get all accounts",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On("GetAll").Return([]models.Account{
					{AccountNumber: "123", AccountHolder: "John"},
				}, nil)
			},
			expected: []models.Account{
				{AccountNumber: "123", AccountHolder: "John"},
			},
		},
		{
			desc:    "ERROR: Failed get all",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On("GetAll").Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &Mocker{repo: mocks.NewIAccountRepo(t)}
			tC.mockSetup(m)

			svc := NewAccountService(m.repo, m.bankRepo)

			result, err := svc.GetAllAccount()

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tC.expected, result)
			}

			m.repo.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAccountByAccNumber(t *testing.T) {
	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
		expected  *models.Account
	}{
		{
			desc:    "SUCCESS",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On("GetByAccountNumber", mock.Anything).
					Return(&models.Account{AccountNumber: "123"}, nil)
			},
			expected: &models.Account{AccountNumber: "123"},
		},
		{
			desc:    "ERROR",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On("GetByAccountNumber", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &Mocker{repo: mocks.NewIAccountRepo(t)}
			tC.mockSetup(m)

			svc := NewAccountService(m.repo, m.bankRepo)

			result, err := svc.GetAccountByAccNumber(t.Context(), "123")

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tC.expected, result)
			}

			m.repo.AssertExpectations(t)
		})
	}
}

func TestAccountService_CreateAcc(t *testing.T) {
	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
	}{
		{
			desc:    "SUCCESS",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On("CreateAcc", mock.Anything).
					Return(&models.Account{AccountNumber: "123"}, nil)
			},
		},
		{
			desc:    "ERROR",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On("CreateAcc", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &Mocker{repo: mocks.NewIAccountRepo(t)}
			tC.mockSetup(m)

			svc := NewAccountService(m.repo, m.bankRepo)

			_, err := svc.CreateAcc(t.Context(), &models.Account{})

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			m.repo.AssertExpectations(t)
		})
	}
}

func TestAccountService_UpdateAcc(t *testing.T) {
	id := uuid.New()

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
	}{
		{
			desc:    "SUCCESS",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On("UpdateAcc", mock.Anything, mock.Anything, mock.Anything).
					Return(&models.Account{}, nil)
			},
		},
		{
			desc:    "ERROR",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On("UpdateAcc", mock.Anything, mock.Anything, mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &Mocker{repo: mocks.NewIAccountRepo(t)}
			tC.mockSetup(m)

			svc := NewAccountService(m.repo, m.bankRepo)

			_, err := svc.UpdateAcc(t.Context(), id, "John", 1000)

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			m.repo.AssertExpectations(t)
		})
	}
}

func TestAccountService_DeleteAcc(t *testing.T) {
	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
	}{
		{
			desc:    "SUCCESS",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.repo.On("DeleteAcc", mock.Anything).
					Return(nil)
			},
		},
		{
			desc:    "ERROR",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.repo.On("DeleteAcc", mock.Anything).
					Return(assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &Mocker{repo: mocks.NewIAccountRepo(t)}
			tC.mockSetup(m)

			svc := NewAccountService(m.repo, m.bankRepo)

			err := svc.DeleteAcc(t.Context(), "123")

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			m.repo.AssertExpectations(t)
		})
	}
}
