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

type BankMocker struct {
	repo *mocks.IBankRepo
}

func TestBankService_GetAll(t *testing.T) {

	testCases := []struct {
		desc      string
		mockSetup func(m *BankMocker)
		wantErr   bool
		expected  []models.Bank
	}{
		{
			desc: "SUCCESS",
			mockSetup: func(m *BankMocker) {
				m.repo.On("GetAll").
					Return([]models.Bank{{Name: "BCA"}}, nil)
			},
			expected: []models.Bank{{Name: "BCA"}},
		},
		{
			desc:    "ERROR",
			wantErr: true,
			mockSetup: func(m *BankMocker) {
				m.repo.On("GetAll").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &BankMocker{repo: mocks.NewIBankRepo(t)}
			tC.mockSetup(m)

			svc := NewBankService(m.repo)

			res, err := svc.GetAllBank()

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tC.expected, res)
			}
		})
	}
}

func TestBankService_GetById(t *testing.T) {
	testCases := []struct {
		desc      string
		mockSetup func(m *BankMocker)
		wantErr   bool
		expected  *models.Bank
	}{
		{
			desc: "SUCCESS",
			mockSetup: func(m *BankMocker) {
				m.repo.On("GetById", "1").
					Return(&models.Bank{Name: "BCA"}, nil)
			},
			expected: &models.Bank{Name: "BCA"},
		},
		{
			desc:    "ERROR",
			wantErr: true,
			mockSetup: func(m *BankMocker) {
				m.repo.On("GetById", "1").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &BankMocker{repo: mocks.NewIBankRepo(t)}
			tC.mockSetup(m)

			svc := NewBankService(m.repo)

			res, err := svc.GetById(t.Context(), "1")

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tC.expected, res)
			}
		})
	}
}

func TestBankService_CreateBank(t *testing.T) {
	testCases := []struct {
		desc      string
		mockSetup func(m *BankMocker)
		wantErr   bool
	}{
		{
			desc: "SUCCESS",
			mockSetup: func(m *BankMocker) {
				m.repo.On("CreateBank", mock.Anything).
					Return(&models.Bank{Name: "BCA"}, nil)
			},
		},
		{
			desc: "ERROR",
			wantErr: true,
			mockSetup: func(m *BankMocker) {
				m.repo.On("CreateBank", mock.Anything).
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &BankMocker{repo: mocks.NewIBankRepo(t)}
			tC.mockSetup(m)

			svc := NewBankService(m.repo)

			_, err := svc.CreateBank(&models.Bank{})

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBankService_UpdateBank(t *testing.T) {
	id := uuid.New()

	testCases := []struct {
		desc      string
		mockSetup func(m *BankMocker)
		wantErr   bool
	}{
		{
			desc: "SUCCESS",
			mockSetup: func(m *BankMocker) {
				m.repo.On("UpdateBank", id, "BCA", "001").
					Return(&models.Bank{Name: "BCA"}, nil)
			},
		},
		{
			desc: "ERROR",
			wantErr: true,
			mockSetup: func(m *BankMocker) {
				m.repo.On("UpdateBank", id, "BCA", "001").
					Return(nil, assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &BankMocker{repo: mocks.NewIBankRepo(t)}
			tC.mockSetup(m)

			svc := NewBankService(m.repo)

			_, err := svc.UpdateBank(id, "BCA", "001")

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBankService_DeleteBank(t *testing.T) {
	testCases := []struct {
		desc      string
		mockSetup func(m *BankMocker)
		wantErr   bool
	}{
		{
			desc: "SUCCESS",
			mockSetup: func(m *BankMocker) {
				m.repo.On("DeleteBank", "1").
					Return(nil)
			},
		},
		{
			desc: "ERROR",
			wantErr: true,
			mockSetup: func(m *BankMocker) {
				m.repo.On("DeleteBank", "1").
					Return(assert.AnError)
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			m := &BankMocker{repo: mocks.NewIBankRepo(t)}
			tC.mockSetup(m)

			svc := NewBankService(m.repo)

			err := svc.DeleteBank("1")

			if tC.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
