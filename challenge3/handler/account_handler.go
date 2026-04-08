package handler

import (
	"challenge3/config"
	"challenge3/dto"
	"challenge3/models"
	"challenge3/service"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	// "github.com/go-chi/chi/v5"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type AccountHandler struct {
	mux                *http.ServeMux
	transactionService *service.TransactionService
	accountService     *service.AccountService
	bankService        *service.BankService
	redis              *redis.Client
}

func NewAccountHandler(mux *http.ServeMux, transactionService *service.TransactionService, accountService *service.AccountService, bankService *service.BankService, redis *redis.Client) *AccountHandler {
	return &AccountHandler{
		mux:                mux,
		accountService:     accountService,
		transactionService: transactionService,
		bankService:        bankService,
		redis:              redis,
	}
}

func (t *AccountHandler) Get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//{service}:{domain}:{identifier}
		key := "account:account:get_all"
		dataCached, err := t.redis.Get(r.Context(), key).Result()
		var accounts = []models.Account{}

		if err == redis.Nil {
			accounts, err = t.accountService.GetAllAccount()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: "error",
					Data:    err.Error(),
				})
				return
			}

			data, _ := json.Marshal(accounts)
			t.redis.Set(r.Context(), key, data, 5*time.Minute)

		} else if err != nil {
			fmt.Printf("unable to GET data from redis. error: %v\n", err)
			accounts, err = t.accountService.GetAllAccount()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: "error",
					Data:    err.Error(),
				})
				return
			}
		} else {
			// cache hit → unmarshal
			if err := json.Unmarshal([]byte(dataCached), &accounts); err != nil {
				fmt.Printf("unable to unmarshal cached data: %v\n", err)
				accounts, err = t.accountService.GetAllAccount()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(dto.BaseResponse{
						Message: "error",
						Data:    err.Error(),
					})
					return
				}
			}
		}

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
			Data:    accounts,
		})
	}
}

func (t *AccountHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context, span := tracer.Start(r.Context(), "account.handler.get-by-id")
		defer span.End()
		id := getIDFromContext(r)
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Print(id)
		// log.Fatal(id)
		account, err := t.accountService.GetAccountById(context, id)

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

var tracer trace.Tracer = otel.Tracer("bank-service")

func (t *AccountHandler) CreateAcc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "AccountHandler")
		traceId := span.SpanContext().TraceID().String()
		defer span.End()
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

		result, err := t.accountService.CreateAcc(ctx, &acc)
		if err != nil {
			if appErr, ok := err.(*dto.ErrorResponse); ok {
				w.WriteHeader(appErr.Code)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: appErr.Message,
				})
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		config.Log.Info("account created",
			zap.String("trace_id", traceId),
		)

		// server.AccountCreatedCounter.Add(r.Context(), 1)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Account Created Successfully",
			Data:    result,
		})
	}
}

func (t *AccountHandler) UpdateAcc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context, span := tracer.Start(r.Context(), "account.handler.update-acc")
		defer span.End()
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

		result, err := t.accountService.UpdateAcc(context, id, req.AccountHolder, req.Balance)
		if err != nil {
			if appErr, ok := err.(*dto.ErrorResponse); ok {
				w.WriteHeader(appErr.Code)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: appErr.Message,
				})
				return
			}
		}

		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Account updated successfully",
			Data:    result,
		})
	}
}

func (t *AccountHandler) DeleteAcc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context, span := tracer.Start(r.Context(), "account.handler.delete-acc")
		defer span.End()
		id := r.URL.Path[len("/accounts/"):]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err := t.accountService.DeleteAcc(context, id)
		if err != nil {
			if appErr, ok := err.(*dto.ErrorResponse); ok {
				w.WriteHeader(appErr.Code)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: appErr.Message,
				})
				return
			}
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
		context, span := tracer.Start(r.Context(), "account.handler.get-history")
		defer span.End()
		idRaw := getIDFromContext(r)
		if idRaw == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if idx := strings.Index(idRaw, "/"); idx != -1 {
			idRaw = idRaw[:idx]
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
		key := "account:transfer:get_history"
		dataCached, err := t.redis.Get(r.Context(), key).Result()
		var transfers = []models.Transaction{}

		if err == redis.Nil {
			transfers, err = t.transactionService.GetHistory(context, id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: "error",
					Data:    err.Error(),
				})
				return
			}

			data, _ := json.Marshal(transfers)
			t.redis.Set(r.Context(), key, data, 5*time.Minute)

		} else if err != nil {
			fmt.Printf("unable to GET data from redis. error: %v\n", err)
			transfers, err = t.transactionService.GetHistory(context, id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: "error",
					Data:    err.Error(),
				})
				return
			}
		} else {
			// cache hit → unmarshal
			if err := json.Unmarshal([]byte(dataCached), &transfers); err != nil {
				fmt.Printf("unable to unmarshal cached data: %v\n", err)
				transfers, err = t.transactionService.GetHistory(context, id)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(w).Encode(dto.BaseResponse{
						Message: "error",
						Data:    err.Error(),
					})
					return
				}
			}
		}
		if err != nil {
			if appErr, ok := err.(*dto.ErrorResponse); ok {
				w.WriteHeader(appErr.Code)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: appErr.Message,
				})
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Transaction History Successfully Fetched",
			Data:    transfers,
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
