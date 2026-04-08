package server

import (
	// "bytes"
	// "context"
	"challenge3/dto"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	// "challenge3/config"
)

func IdempotencyMiddlewareRedis(redis *redis.Client, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ctx := context.Background()

		if r.Method != http.MethodPost && r.Method != http.MethodPut {
			next.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("Idempotency-Key")
		if key == "" {
			next.ServeHTTP(w, r)
			return
		}

		redisKey := "server:idempotency:key:" + key
		
		cached, err := redis.Get(r.Context(), redisKey).Result()
		if err == nil && cached != "" {
			var resp struct {
				Status int
				Body   []byte
			}
			if err := json.Unmarshal([]byte(cached), &resp); err == nil {
				w.WriteHeader(resp.Status)
				w.Write(resp.Body)
				return
			}
		}

		// 2️⃣ record response
		rec := &dto.ResponseRecorder{
			ResponseWriter: w,
			Status:         200,
			Body:           []byte{},
		}

		next.ServeHTTP(rec, r)

		data, _ := json.Marshal(map[string]interface{}{
			"status": rec.Status,
			"body":   rec.Body,
		})

		// TTL 10 menit
		redis.Set(r.Context(), redisKey, data, 10*time.Minute)
	})
}
