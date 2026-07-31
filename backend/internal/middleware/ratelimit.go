package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/murovei88/pw-toolkit/pkg/httputil"
)

// RateLimiter создаёт middleware для ограничения запросов
// maxRequests — максимум запросов за period с одного IP
func RateLimiter(rdb *redis.Client, maxRequests int, period time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)
			key := fmt.Sprintf("ratelimit:%s:%s", r.URL.Path, ip)

			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()

			// Инкрементируем счётчик
			count, err := rdb.Incr(ctx, key).Result()
			if err != nil {
				// Если Redis упал — пропускаем (fail open)
				// В production лучше fail closed
				next.ServeHTTP(w, r)
				return
			}

			// Устанавливаем TTL только на первый запрос
			if count == 1 {
				rdb.Expire(ctx, key, period)
			}

			// Проверяем лимит
			if count > int64(maxRequests) {
				httputil.JSON(w, http.StatusTooManyRequests, httputil.APIResponse{
					Status:  429,
					Success: false,
					Message: fmt.Sprintf("rate limit exceeded: %d requests per %v", maxRequests, period),
				})
				return
			}

			// Добавляем заголовки с информацией о лимите
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", maxRequests))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", int64(maxRequests)-count))

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP извлекает реальный IP клиента
func getClientIP(r *http.Request) string {
	// Проверяем X-Forwarded-For (если за nginx)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	
	// Проверяем X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback: RemoteAddr
	addr := r.RemoteAddr
	if idx := strings.LastIndex(addr, ":"); idx != -1 {
		return addr[:idx]
	}
	return addr
}
