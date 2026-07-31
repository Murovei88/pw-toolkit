package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/murovei88/pw-toolkit/internal/domain"
	"github.com/murovei88/pw-toolkit/internal/service"
	"github.com/murovei88/pw-toolkit/pkg/httputil"
)

type BuildHandler struct {
	buildService *service.BuildService
}

func NewBuildHandler(buildService *service.BuildService) *BuildHandler {
	return &BuildHandler{buildService: buildService}
}

// CreateBuildRequest — DTO для создания билда
type CreateBuildRequest struct {
	Name        string              `json:"name"`
	ClassID     int                 `json:"class_id"`
	Level       int                 `json:"level"`
	Equipment   domain.Equipment    `json:"equipment"`
	Cards       []int               `json:"cards"`
	Books       []int               `json:"books"`
	GenieID     *int                `json:"genie_id"`
	PanguSouls  domain.PanguSouls   `json:"pangu_souls"`
	StarDisks   domain.StarDisks    `json:"star_disks"`
	Titles      []int               `json:"titles"`
}

// CreateBuild POST /api/v1/builds
func (h *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var req CreateBuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.BadRequest(w, "invalid JSON body")
		return
	}
	defer r.Body.Close()

	build := &domain.Build{
		Name:        req.Name,
		ClassID:     req.ClassID,
		Level:       req.Level,
		Equipment:   req.Equipment,
		Cards:       req.Cards,
		Books:       req.Books,
		GenieID:     req.GenieID,
		PanguSouls:  req.PanguSouls,
		StarDisks:   req.StarDisks,
		Titles:      req.Titles,
	}

	if err := h.buildService.CreateBuild(r.Context(), build); err != nil {
		if errors.Is(err, service.ErrInvalidBuild) {
			httputil.BadRequest(w, err.Error())
			return
		}
		httputil.InternalError(w, "failed to create build")
		return
	}

	httputil.JSON(w, http.StatusCreated, httputil.APIResponse{
		Status:  201,
		Success: true,
		Message: "Build created",
		Data: map[string]interface{}{
			"id":  build.ID,
			"url": "/builds/" + build.ID,
		},
	})
}

// GetBuild GET /api/v1/builds/{id}
func (h *BuildHandler) GetBuild(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL path
	// Используем стандартный паттерн Go 1.22+ path values
	id := r.PathValue("id")
	if id == "" {
		// Fallback: парсим вручную
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) >= 4 {
			id = parts[len(parts)-1]
		}
	}

	if id == "" {
		httputil.BadRequest(w, "build ID is required")
		return
	}

	build, err := h.buildService.GetBuild(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrBuildNotFound) {
			httputil.NotFound(w, "build not found")
			return
		}
		httputil.InternalError(w, "failed to get build")
		return
	}

	httputil.Success(w, "OK", build)
}
