package handler

import (
	"net/http"
	"runtime"

	"github.com/murovei88/pw-toolkit/internal/config"
	"github.com/murovei88/pw-toolkit/pkg/httputil"
)

type StatusData struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	Commit    string `json:"commit"`
	Branch    string `json:"branch"`
	GoVersion string `json:"goVersion"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
}

func StatusHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httputil.Success(w, "OK", StatusData{
			Name:      "pw-toolkit",
			Version:   cfg.Version,
			BuildDate: cfg.BuildDate,
			Commit:    cfg.Commit,
			Branch:    cfg.Branch,
			GoVersion: runtime.Version(),
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
		})
	}
}
