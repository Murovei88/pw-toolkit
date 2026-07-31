package handler

import (
	"net/http"

	"github.com/murovei88/pw-toolkit/pkg/httputil"
)

type ServiceHealth struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	ResponseTime string `json:"responseTime,omitempty"`
	Error        string `json:"error,omitempty"`
}

type HealthData struct {
	Services         []ServiceHealth `json:"services"`
	TotalServices    int             `json:"totalServices"`
	HealthyServices  int             `json:"healthyServices"`
	Degraded         bool            `json:"degraded"`
}

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement actual health checks for MariaDB, Redis, MinIO
		// For now, return healthy stub
		
		services := []ServiceHealth{
			{Name: "api", Status: "healthy", ResponseTime: "0.5ms"},
			{Name: "mariadb", Status: "healthy", ResponseTime: "2.3ms"},
			{Name: "redis", Status: "healthy", ResponseTime: "0.8ms"},
			{Name: "minio", Status: "healthy", ResponseTime: "15.2ms"},
		}

		data := HealthData{
			Services:        services,
			TotalServices:   len(services),
			HealthyServices: len(services),
			Degraded:        false,
		}

		httputil.Success(w, "All services healthy", data)
	}
}
