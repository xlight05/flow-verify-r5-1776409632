package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/asdlc/todo-api/internal/models"
	"github.com/asdlc/todo-api/internal/store"
	"github.com/google/uuid"
)

const (
	codeValidation = "VALIDATION_ERROR"
	codeInvalidJSON = "INVALID_JSON"
	codeNotFound   = "NOT_FOUND"
	codeInternal   = "INTERNAL_ERROR"
)

type Handler struct {
	Store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{Store: s}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/todos", h.todosCollection)
	mux.HandleFunc("/todos/", h.todosItem)
	return mux
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, codeValidation, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) todosCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		todos := h.Store.List()
		writeJSON(w, http.StatusOK, todos)
	case http.MethodPost:
		h.createTodo(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, codeValidation, "method not allowed")
	}
}

func (h *Handler) todosItem(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todos/")
	id = strings.Trim(id, "/")
	if id == "" || strings.Contains(id, "/") {
		writeError(w, http.StatusNotFound, codeNotFound, "todo not found")
		return
	}
	switch r.Method {
	case http.MethodDelete:
		if err := h.Store.Delete(id); err != nil {
			if errors.Is(err, store.ErrNotFound) {
				writeError(w, http.StatusNotFound, codeNotFound, "todo not found")
				return
			}
			log.Printf("delete error: %v", err)
			writeError(w, http.StatusInternalServerError, codeInternal, "internal server error")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeError(w, http.StatusMethodNotAllowed, codeValidation, "method not allowed")
	}
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "invalid request body")
		return
	}
	var req models.CreateTodoRequest
	dec := json.NewDecoder(strings.NewReader(string(body)))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, codeInvalidJSON, "invalid json")
		return
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, http.StatusBadRequest, codeValidation, "title is required")
		return
	}
	if len(title) > 255 {
		writeError(w, http.StatusBadRequest, codeValidation, "title must be at most 255 characters")
		return
	}
	completed := false
	if req.Completed != nil {
		completed = *req.Completed
	}
	todo := models.Todo{
		ID:        uuid.NewString(),
		Title:     title,
		Completed: completed,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := h.Store.Create(todo); err != nil {
		log.Printf("create error: %v", err)
		writeError(w, http.StatusInternalServerError, codeInternal, "internal server error")
		return
	}
	writeJSON(w, http.StatusCreated, todo)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, msg string) {
	writeJSON(w, status, models.ErrorResponse{Error: msg, Code: code})
}
