package domain

import "time"

type Employee struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Position  string    `json:"position"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Employee) Validate() error {
	if e.Name == "" {
		return ErrValidation{"name required"}
	}
	if e.Age <= 0 {
		return ErrValidation{"age must be > 0"}
	}
	return nil
}

type ErrValidation struct{ Msg string }

func (e ErrValidation) Error() string { return e.Msg }
