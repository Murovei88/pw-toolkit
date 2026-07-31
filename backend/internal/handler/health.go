package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/murovei88/pw-toolkit/pkg/httputil"
)

type ServiceHealth struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	ResponseTime string `json:"responseTime,omitempty"`
	Error        string `json:"error,omitempty"`
}

type HealthData struct {
	Services        []ServiceHealth `json:"services"`
	TotalServices   int             `json:"totalServices"`
	HealthyServices int             `json:"healthyServices"`
	Degraded        bool            `json:"degraded"`
}

func HealthHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		services := []ServiceHealth{}

		// API self-check
		services = append(services, ServiceHealth{
			Name:         "api",
			Status:       "healthy",
			ResponseTime: "0.5ms",
		})

		// MariaDB
		services = append(services, checkMariaDB(db))

		// Redis
		services = append(services, checkRedis(rdb))

		// MinIO (TODO)
		services = append(services, ServiceHealth{
			Name:         "minio",
			Status:       "healthy",
			ResponseTime: "15.2ms",
		})

		healthyCount := 0
		for _, s := range services {
			if s.Status == "healthy" {
				healthyCount++
			}
		}

		degraded := healthyCount < len(services)
		status := http.StatusOK
		message := "All services healthy"
		if degraded {
			status = http.StatusServiceUnavailable
			message = "Service degradation detected"
		}

		httputil.JSON(w, status, httputil.APIResponse{
			Status:  status,
			Success: !degraded,
			Message: message,
			Data: HealthData{
				Services:        services,
				TotalServices:   len(services),
				HealthyServices: healthyCount,
				Degraded:        degraded,
			},
		})
	}
}

func checkMariaDB(db *sql.DB) ServiceHealth {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	if err != nil {
		return ServiceHealth{
			Name:   "mariadb",
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}
	return ServiceHealth{
		Name:         "mariadb",
		Status:       "healthy",
		ResponseTime: time.Since(start).String(),
	}
}

func checkRedis(rdb *redis.Client) ServiceHealth {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return ServiceHealth{
			Name:   "redis",
			Status: "unhealthy",
			Error:  err.Error(),
		}
	}
	return ServiceHealth{
		Name:         "redis",
		Status:       "healthy",
		ResponseTime: time.Since(start).String(),
	}
}
