package handler

import (
	"challenge3/dto"
	"challenge3/models"
	"challenge3/service"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	// "github.com/go-chi/chi/v5"
	"fmt"
	"github.com/google/uuid"
	"io"
)

type AccountHandler struct {
	mux                *http.ServeMux
	transactionService *service.TransactionService
	accountService     *service.AccountService
	bankService        *service.BankService
}

func NewAccountHandler(mux *http.ServeMux, transactionService *service.TransactionService, accountService *service.AccountService, bankService *service.BankService) *AccountHandler {
	return &AccountHandler{
		mux:                mux,
		accountService:     accountService,
		transactionService: transactionService,
		bankService:        bankService,
	}
}

func (t *AccountHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := t.accountService.GetAllAccount()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "succes",
			Data:    account,
		})
	}
}

func (t *AccountHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getIDFromContext(r)
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Print(id)
		// log.Fatal(id)
		account, err := t.accountService.GetAccountById(id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "succes",
			Data:    account,
		})
	}
}

func (t *AccountHandler) CreateAcc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var acc models.Account

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

		if err := json.Unmarshal(body, &acc); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Invalid JSON",
			})
			return
		}

		_, err = t.bankService.GetById(acc.BankID.String())
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Bank ID doesnt exists",
			})
			return
		}

		existing, err := t.accountService.GetAccountByAccNumber(acc.AccountNumber)
		if err == nil && existing != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Account with this number already exists",
			})
			return
		}

		if acc.Balance <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Balance cant be lower than 0",
				// Data:    err.Error(),
			})
			return
		}

		result, err := t.accountService.CreateAcc(&acc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Account Created Successfully",
			Data:    result,
		})
	}
}

func (t *AccountHandler) UpdateAcc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// r.PathValue(r)
		var req struct {
			AccountHolder string `json:"account_holder"`
			Balance       int    `json:"balance"`
		}

		idRaw := r.URL.Path[len("/accounts/"):]
		if idRaw == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id, err := uuid.Parse(idRaw)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		body, err := io.ReadAll(r.Body)
		fmt.Print(body)
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

		existing, err := t.accountService.GetAccountById(idRaw)
		if err != nil || existing == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if existing.AccountHolder != req.AccountHolder {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Account Holder not Match with the Data",
				// Data:    err.Error(),
			})
			return
		}

		if req.Balance <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Balance cant be lower than 0",
				// Data:    err.Error(),
			})
			return
		}

		result, err := t.accountService.UpdateAcc(id, req.AccountHolder, req.Balance)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Account updated successfully",
			Data:    result,
		})
	}
}

func (t *AccountHandler) DeleteAcc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/accounts/"):]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := t.accountService.DeleteAcc(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Account with this id doesnt exists",
			})
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Account Deleted Successfully",
			// Data:    result,
		})
	}
}

func (t *AccountHandler) GetHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := getIDFromContext(r)
		if idRaw == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if idx := strings.Index(idRaw, "/"); idx != -1 {
			idRaw = idRaw[:idx]
		}

		fmt.Print(idRaw)
		// log.Fatal(idRaw)

		id, err := uuid.Parse(idRaw)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		_, err = t.accountService.GetAccountById(idRaw)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Account doesnt exist",
				Data:    err.Error(),
			})
			return
		}

		histories, err := t.transactionService.GetHistory(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "error",
				Data:    err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Transaction History Successfully Fetched",
			Data:    histories,
		})
	}
}

type contextKey string

const idKey contextKey = "id"

func setIDToContext(r *http.Request, id string) *http.Request {
	ctx := context.WithValue(r.Context(), idKey, id)
	return r.WithContext(ctx)
}

func getIDFromContext(r *http.Request) string {
	id, _ := r.Context().Value(idKey).(string)
	return id
}

func (t *AccountHandler) AccountRouter() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/accounts/")
		parts := strings.Split(path, "/")

		// contoh:
		// /accounts/123 -> ["123"]
		// /accounts/123/transfer -> ["123", "transfer"]

		if len(parts) == 1 {
			// GET /accounts/{id}
			r = setIDToContext(r, parts[0]) // optional helper
			t.GetById()(w, r)
			return
		}

		if len(parts) == 2 && parts[1] == "transfer" {
			// GET /accounts/{id}/transfer
			r = setIDToContext(r, parts[0])
			t.GetHistory()(w, r)
			return
		}

		http.NotFound(w, r)
	}
}
