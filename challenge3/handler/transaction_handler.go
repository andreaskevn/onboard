package handler

import (
	"challenge3/dto"
	// "challenge3/models"
	"challenge3/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TransactionHandler struct {
	mux                *http.ServeMux
	transactionService *service.TransactionService
	accountService     *service.AccountService
	redis              *redis.Client
}

func NewTransctionHandler(mux *http.ServeMux, transactionService *service.TransactionService, accountService *service.AccountService, redis *redis.Client) *TransactionHandler {
	return &TransactionHandler{
		mux:                mux,
		transactionService: transactionService,
		accountService:     accountService,
		redis:              redis,
	}
}

func (t *TransactionHandler) Transfer() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context, span := tracer.Start(r.Context(), "TransactionHandler")
		defer span.End()
		var req struct {
			AccountFrom uuid.UUID `json:"from_account_id"`
			AccountTo   uuid.UUID `json:"to_account_id"`
			Amount      int       `json:"amount"`
			AdminFee    int       `json:"admin_fee"`
		}

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{
				Message: "Idempotency Key is missing",
			})
			return
		}

		redisKey := "transaction:transfer:id:" + key
		cached, err := t.redis.Get(r.Context(), redisKey).Result()
		if err == nil && cached != "" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(cached))
			return
		} else if err != nil && err != redis.Nil {
			fmt.Printf("redis error: %v\n", err)
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

		result, err := t.transactionService.Transfer(context, req.AccountFrom, req.AccountTo, req.Amount, req.AdminFee)
		if err != nil {
			if appErr, ok := err.(*dto.ErrorResponse); ok {
				w.WriteHeader(appErr.Code)
				json.NewEncoder(w).Encode(dto.BaseResponse{
					Message: appErr.Message,
				})
				return
			}
			// fmt.Print("error: ", err)
		}
		// res, _ := json.Marshal(res)
		t.redis.Set(r.Context(), redisKey, result, 10*time.Minute)

		result.AdminFee = req.AdminFee - (req.AdminFee * 2)

		json.NewEncoder(w).Encode(dto.BaseResponse{
			Message: "Transfer successfully sent",
			Data:    result,
		})
	}
}
