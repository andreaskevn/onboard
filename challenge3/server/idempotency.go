package server

import (
	"errors"
	"net/http"

	"challenge3/dto"
	"challenge3/models"

	"fmt"
	"gorm.io/gorm"
)

func IdempotencyMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Print(key)

		var record models.IdempotencyKey
		err := db.First(&record, "id = ?", key).Error

		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(record.Status)
			w.Write([]byte(record.Response))
			return
		}

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		rec := &dto.ResponseRecorder{
			ResponseWriter: w,
			Status:         200,
			Body:           []byte{},
		}

		next.ServeHTTP(rec, r)

		// simpan (ignore error duplicate)
		record = models.IdempotencyKey{
			ID:       key,
			Response: string(rec.Body),
			Status:   rec.Status,
		}

		_ = db.Create(&record).Error
	})
}
