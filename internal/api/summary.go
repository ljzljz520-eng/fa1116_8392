package api

import (
	"gestureparticles/internal/model"
	"gestureparticles/internal/registry"
	"net/http"
)

func SummaryHandler(service *registry.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		snapshot, err := service.Snapshot(r.URL.Query().Get("classroom"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, snapshot)
	}
}

func statusName(status model.RecordStatus) string {
	if status == "" {
		return string(model.StatusDraft)
	}
	return string(status)
}
