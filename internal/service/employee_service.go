package service

import (
	"context"
	"practice-go-crud/internal/domain"
	"practice-go-crud/internal/repository"
	"time"
)

type EmployeeService interface {
	List(ctx context.Context) ([]domain.Employee, error)
	Get(ctx context.Context, id int) (*domain.Employee, error)
	Create(ctx context.Context, e *domain.Employee) error
	Update(ctx context.Context, e *domain.Employee) error
	Delete(ctx context.Context, id int) error
}

type employeeService struct{ repo repository.EmployeeRepository }

func NewEmployeeService(r repository.EmployeeRepository) EmployeeService {
	return &employeeService{repo: r}
}

func (s *employeeService) List(ctx context.Context) ([]domain.Employee, error) {
	return s.repo.FindAll(ctx)
}
func (s *employeeService) Get(ctx context.Context, id int) (*domain.Employee, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *employeeService) Create(ctx context.Context, e *domain.Employee) error {
	if err := e.Validate(); err != nil {
		return err
	}
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
	return s.repo.Create(ctx, e)
}

func (s *employeeService) Update(ctx context.Context, e *domain.Employee) error {
	if err := e.Validate(); err != nil {
		return err
	}
	e.UpdatedAt = time.Now()
	return s.repo.Update(ctx, e)
}

func (s *employeeService) Delete(ctx context.Context, id int) error { return s.repo.Delete(ctx, id) }
