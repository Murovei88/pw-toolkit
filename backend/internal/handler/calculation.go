package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/murovei88/pw-toolkit/internal/calc"
	"github.com/murovei88/pw-toolkit/internal/service"
	"github.com/murovei88/pw-toolkit/pkg/httputil"
)

type CalculationHandler struct {
	buildService *service.BuildService
	logger       *slog.Logger
}

func NewCalculationHandler(buildService *service.BuildService, logger *slog.Logger) *CalculationHandler {
	return &CalculationHandler{
		buildService: buildService,
		logger:       logger,
	}
}

// CalculatePreview POST /api/v1/calculate
func (h *CalculationHandler) CalculatePreview(w http.ResponseWriter, r *http.Request) {
	var req calc.PreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.BadRequest(w, "invalid JSON body")
		return
	}
	defer r.Body.Close()

	stats, err := h.buildService.CalculatePreview(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to calculate preview", "error", err)
		httputil.InternalError(w, "failed to calculate stats")
		return
	}

	httputil.Success(w, "OK", stats)
}
