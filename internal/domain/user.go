package domain

import (
	"errors"
	"regexp"
	"time"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type AdminRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

var (
	ErrInvalidEmail     = errors.New("email format invalid")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrNameRequired     = errors.New("name is required")
	ErrEmailRequired    = errors.New("email is required")
	ErrPasswordRequired = errors.New("password is required")
)

func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if r.Name == "" {
		return ErrNameRequired
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return ErrInvalidEmail
	}

	if len(r.Password) < 6 {
		return ErrPasswordTooShort
	}

	return nil
}

func (r *AdminRegisterRequest) Validate() error {
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	if r.Name == "" {
		return ErrNameRequired
	}
	if r.Role != "user" && r.Role != "admin" {
		return errors.New("role must be either be 'user' or 'admin'")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return ErrInvalidEmail
	}

	if len(r.Password) < 6 {
		return ErrPasswordTooShort
	}

	return nil
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return ErrEmailRequired
	}
	if r.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}
