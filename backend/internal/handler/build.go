package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/murovei88/pw-toolkit/internal/domain"
	"github.com/murovei88/pw-toolkit/internal/service"
	"github.com/murovei88/pw-toolkit/pkg/httputil"
)

type BuildHandler struct {
	buildService *service.BuildService
	logger       *slog.Logger
}

func NewBuildHandler(buildService *service.BuildService, logger *slog.Logger) *BuildHandler {
	return &BuildHandler{
		buildService: buildService,
		logger:       logger,
	}
}

type CreateBuildRequest struct {
	Name       string            `json:"name"`
	ClassID    int               `json:"class_id"`
	Level      int               `json:"level"`
	Equipment  domain.Equipment  `json:"equipment"`
	Cards      []int             `json:"cards"`
	Books      []int             `json:"books"`
	GenieID    *int              `json:"genie_id"`
	PanguSouls domain.PanguSouls `json:"pangu_souls"`
	StarDisks  domain.StarDisks  `json:"star_disks"`
	Titles     []int             `json:"titles"`
}

func (h *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var req CreateBuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.BadRequest(w, "invalid JSON body")
		return
	}
	defer r.Body.Close()

	build := &domain.Build{
		Name:       req.Name,
		ClassID:    req.ClassID,
		Level:      req.Level,
		Equipment:  req.Equipment,
		Cards:      req.Cards,
		Books:      req.Books,
		GenieID:    req.GenieID,
		PanguSouls: req.PanguSouls,
		StarDisks:  req.StarDisks,
		Titles:     req.Titles,
	}

	if err := h.buildService.CreateBuild(r.Context(), build); err != nil {
		h.logger.Error("Failed to create build",
			"error", err,
			"class_id", req.ClassID,
			"level", req.Level,
		)

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

func (h *BuildHandler) GetBuild(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
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
		h.logger.Error("Failed to get build", "id", id, "error", err)
		httputil.InternalError(w, "failed to get build")
		return
	}

	httputil.Success(w, "OK", build)
}
