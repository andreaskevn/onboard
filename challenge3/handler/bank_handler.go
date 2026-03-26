package handler

import (
	"challenge3/dto"
	"challenge3/models"
	"challenge3/service"
	"encoding/json"
	"net/http"

	// "github.com/go-chi/chi/v5"
	"fmt"
	"github.com/google/uuid"
	"io"
)

type BankHandler struct {
	mux         *http.ServeMux
	bankService *service.BankService
}

func NewBankHandler(mux *http.ServeMux, bankService *service.BankService) *BankHandler {
	return &BankHandler{
		mux:         mux,
		bankService: bankService,
	}
}

func (t *BankHandler) GetBank() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := t.bankService.GetAllBank()
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

func (t *BankHandler) GetBankById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := getIDFromContext(r)
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Print(id)
		// log.Fatal(id)
		account, err := t.bankService.GetById(id)

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

func (t *BankHandler) CreateAcc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var bank models.Bank

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

		if err := json.Unmarshal(body, &bank); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Invalid JSON",
			})
			return
		}

		bankCodeExist, err := t.bankService.GetByCode(bank.Code)
		fmt.Print(bank.Code)
		fmt.Print(bank.Name)
		if err == nil && bankCodeExist != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Bank with this Code already exists",
			})
			fmt.Println("bankCodeExist: %w", bankCodeExist)
			return
		}

		bankNameExist, err := t.bankService.GetByName(bank.Name)
		if err == nil && bankNameExist != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Bank with this Name already exists",
			})
			return
		}

		result, err := t.bankService.CreateBank(&bank)
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

func (t *BankHandler) UpdateBank() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			BankName string `json:"bank_name"`
			BankCode string `json:"bank_code"`
		}

		idRaw := r.URL.Path[len("/banks/"):]
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

		existing, err := t.bankService.GetById(idRaw)
		if err != nil || existing == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		bankCodeExist, err := t.bankService.GetByCode(req.BankCode)
		if err == nil && bankCodeExist != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Bank with this Code already exists",
			})
			return
		}

		bankNameExist, err := t.bankService.GetByCode(req.BankName)
		if err == nil && bankNameExist != nil {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Bank with this Name already exists",
			})
			return
		}

		result, err := t.bankService.UpdateBank(id, req.BankName, req.BankCode)
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

func (t *BankHandler) DeleteBank() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/banks/"):]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := t.bankService.DeleteBank(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Bank with this id doesnt exists",
			})
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Bank Deleted Successfully",
			// Data:    result,
		})
	}
}
