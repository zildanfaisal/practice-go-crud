package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"practice-go-crud/internal/domain"
	"practice-go-crud/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.authService.Register(r.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrEmailRequired, domain.ErrPasswordRequired, domain.ErrNameRequired:
			writeError(w, http.StatusBadRequest, err.Error())
		case domain.ErrInvalidEmail, domain.ErrPasswordTooShort:
			writeError(w, http.StatusBadRequest, err.Error())
		case service.ErrEmailAlreadyExists:
			writeError(w, http.StatusConflict, "Email already registered")
		default:
			writeError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (h *AuthHandler) AdminRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req domain.AdminRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	user, err := h.authService.AdminRegister(r.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrEmailRequired, domain.ErrPasswordRequired, domain.ErrNameRequired:
			writeError(w, http.StatusBadRequest, err.Error())
		case domain.ErrInvalidEmail, domain.ErrPasswordTooShort:
			writeError(w, http.StatusBadRequest, err.Error())
		case service.ErrEmailAlreadyExists:
			writeError(w, http.StatusConflict, "Email already registered")
		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Admin user registered successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	loginResponse, err := h.authService.Login(r.Context(), &req)
	if err != nil {
		switch err {
		case domain.ErrEmailRequired, domain.ErrPasswordRequired:
			writeError(w, http.StatusBadRequest, err.Error())
		case service.ErrInvalidCredentials:
			writeError(w, http.StatusUnauthorized, "Invalid email or password")
		default:
			writeError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Login successful",
		"token":   loginResponse.Token,
		"user":    loginResponse.User,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := h.extractTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "Authorization token required")
		return
	}

	user, err := h.authService.ValidateToken(token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User profile retrieved successfully",
		"user":    user,
	})
}

func (h *AuthHandler) extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := h.extractTokenFromHeader(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "Authorization token required")
		return
	}

	err := h.authService.Logout(r.Context(), token)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to logout: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout successful",
	})
}
