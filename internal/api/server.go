package api

import (
	"encoding/json"
	"gestureparticles/internal/flow021"
	"gestureparticles/internal/importer"
	"gestureparticles/internal/lesson"
	"gestureparticles/internal/model"
	"gestureparticles/internal/particle"
	"gestureparticles/internal/registry"
	"net/http"
	"strings"
)

type Server struct {
	service *registry.Service
	mux     *http.ServeMux
}

func New(service *registry.Service) http.Handler {
	s := &Server{service: service, mux: http.NewServeMux()}
	s.routes()
	return logging(s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/health", s.health)
	s.mux.HandleFunc("/records", s.records)
	s.mux.HandleFunc("/records/", s.record)
	s.mux.HandleFunc("/observations", s.observations)
	s.mux.HandleFunc("/workflows", s.workflows)
	s.mux.HandleFunc("/import", s.importRecords)
	s.mux.HandleFunc("/lessons/schedule", s.lessonSchedule)
	s.mux.HandleFunc("/frames", s.frames)
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) records(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		values, err := s.service.Search(r.URL.Query().Get("classroom"), r.URL.Query().Get("q"), model.RecordStatus(r.URL.Query().Get("status")))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, values)
	case http.MethodPost:
		var input model.Record
		if err := decode(r, &input); err != nil {
			writeError(w, err)
			return
		}
		value, err := s.service.Register(input)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, value)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) record(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/records/")
	if id == "" {
		writeError(w, errBadRequest("record id required"))
		return
	}
	switch r.Method {
	case http.MethodGet:
		value, err := s.service.Get(id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, value)
	case http.MethodPatch:
		var input model.Record
		if err := decode(r, &input); err != nil {
			writeError(w, err)
			return
		}
		value, err := s.service.Update(id, input, r.Header.Get("X-Actor"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, value)
	case http.MethodPost:
		if r.URL.Query().Get("action") == "review" {
			value, err := s.service.Review(id, r.Header.Get("X-Actor"), r.URL.Query().Get("decision"), r.URL.Query().Get("reason"))
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, value)
			return
		}
		if r.URL.Query().Get("action") == "archive" {
			value, err := s.service.Archive(id, r.Header.Get("X-Actor"))
			if err != nil {
				writeError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, value)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) observations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var input flow021.Observation
	if err := decode(r, &input); err != nil {
		writeError(w, err)
		return
	}
	value, err := s.service.AddObservation(input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, value)
}

func (s *Server) workflows(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var input model.Workflow
	if err := decode(r, &input); err != nil {
		writeError(w, err)
		return
	}
	value, err := s.service.StartWorkflow(input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, value)
}

func (s *Server) importRecords(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var rows []importer.Row
	if err := decode(r, &rows); err != nil {
		writeError(w, err)
		return
	}
	report, err := s.service.Import(rows, r.Header.Get("X-Actor"))
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, report)
}

func (s *Server) lessonSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var plan lesson.Plan
	if err := decode(r, &plan); err != nil {
		writeError(w, err)
		return
	}
	slots, err := s.service.PrepareLesson(plan)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, slots)
}

func (s *Server) frames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var observations []flow021.Observation
	if err := decode(r, &observations); err != nil {
		writeError(w, err)
		return
	}
	sequence := particle.BuildSequence(observations)
	writeJSON(w, http.StatusOK, sequence)
}

func decode(r *http.Request, target any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

type badRequest string

func errBadRequest(message string) error { return badRequest(message) }

func (e badRequest) Error() string { return string(e) }

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") == "" {
			w.Header().Set("X-Request-ID", "deterministic")
		}
		next.ServeHTTP(w, r)
	})
}
