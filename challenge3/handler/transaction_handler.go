package handler

import (
	"challenge3/dto"
	// "challenge3/models"
	"challenge3/service"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type TransactionHandler struct {
	mux                *http.ServeMux
	transactionService *service.TransactionService
	accountService     *service.AccountService
}

func NewTransctionHandler(mux *http.ServeMux, transactionService *service.TransactionService, accountService *service.AccountService) *TransactionHandler {
	return &TransactionHandler{
		mux:                mux,
		transactionService: transactionService,
		accountService:     accountService,
	}
}

func (t *TransactionHandler) Transfer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			AccountFrom uuid.UUID `json:"from_account_id"`
			AccountTo   uuid.UUID `json:"to_account_id"`
			Amount      int       `json:"amount"`
			AdminFee    int       `json:"admin_fee"`
		}

		body, err := io.ReadAll(r.Body)
		// fmt.Print(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(body) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Request body is empty",
			})
			return
		}

		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Invalid JSON",
			})
			return
		}

		if req.AccountFrom == req.AccountTo {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Account From and Destination cant be same",
				// Data:    err.Error(),
			})
			return
		}

		accountFromExist, err := t.accountService.GetAccountById(req.AccountFrom.String())
		if err != nil || accountFromExist == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Account From doesnt exist",
				// Data:    err.Error(),
			})
			return
		}

		accountToExist, err := t.accountService.GetAccountById(req.AccountTo.String())
		if err != nil || accountToExist == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Account Destination doesnt exist",
				// Data:    err.Error(),
			})
			return
		}

		if req.Amount <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Amount cant be lower than 0",
				// Data:    err.Error(),
			})
			return
		}

		if accountFromExist.Balance < (req.Amount + req.AdminFee) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Your balance is insufficient",
				// Data:    err.Error(),
			})
			return
		}

		if accountFromExist.BankID != accountToExist.BankID {
			if req.AdminFee <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: "Admin Fee must included",
					// Data:    accountToExist,
				})
				return
			}
		}
		fmt.Print(&accountFromExist)
		fmt.Print(&accountToExist)

		result, err := t.transactionService.Transfer(req.AccountFrom, req.AccountTo, req.Amount, req.AdminFee)
		result.AdminFee = req.AdminFee - (req.AdminFee * 2)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		reduceBalance := accountFromExist.Balance - (req.Amount + req.AdminFee)
		_, err = t.accountService.UpdateAcc(accountFromExist.ID, accountFromExist.AccountHolder, reduceBalance)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Transfer successfully sent",
			Data:    result,
		})
	}
}
