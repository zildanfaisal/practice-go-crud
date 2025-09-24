package handler

import (
	"encoding/json"
	"net/http"
	"practice-go-crud/internal/domain"
	"practice-go-crud/internal/service"
	"strconv"
)

type EmployeeHandler struct{ Svc service.EmployeeService }

func NewEmployeeHandler(s service.EmployeeService) *EmployeeHandler { return &EmployeeHandler{Svc: s} }

func (h *EmployeeHandler) List(w http.ResponseWriter, r *http.Request) {
	ems, err := h.Svc.List(r.Context())
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	writeJSON(w, 200, ems)
}
func (h *EmployeeHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	e, err := h.Svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, 500, err.Error())
		return
	}
	if e == nil {
		writeError(w, 404, "not found")
		return
	}
	writeJSON(w, 200, e)
}
func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var e domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, 400, "invalid json")
		return
	}
	if err := h.Svc.Create(r.Context(), &e); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	writeJSON(w, 201, e)
}
func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	var e domain.Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, 400, "invalid json")
		return
	}
	e.ID = id
	if err := h.Svc.Update(r.Context(), &e); err != nil {
		writeError(w, 400, err.Error())
		return
	}
	writeJSON(w, 200, e)
}
func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.Svc.Delete(r.Context(), id); err != nil {
		writeError(w, 500, err.Error())
		return
	}
	w.WriteHeader(204)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
