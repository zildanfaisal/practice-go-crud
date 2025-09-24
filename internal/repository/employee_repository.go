package repository

import (
	"context"
	"database/sql"
	"errors"
	"practice-go-crud/internal/domain"
)

type EmployeeRepository interface {
	FindAll(ctx context.Context) ([]domain.Employee, error)
	FindByID(ctx context.Context, id int) (*domain.Employee, error)
	Create(ctx context.Context, e *domain.Employee) error
	Update(ctx context.Context, e *domain.Employee) error
	Delete(ctx context.Context, id int) error
}

type mysqlEmployeeRepo struct{ db *sql.DB }

func NewMySQLEmployeeRepo(db *sql.DB) EmployeeRepository { return &mysqlEmployeeRepo{db: db} }

func (r *mysqlEmployeeRepo) FindAll(ctx context.Context) ([]domain.Employee, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, age, position, address, created_at, updated_at FROM employee")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	emps := []domain.Employee{}
	for rows.Next() {
		var e domain.Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Age, &e.Position, &e.Address, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		emps = append(emps, e)
	}
	return emps, nil
}

func (r *mysqlEmployeeRepo) FindByID(ctx context.Context, id int) (*domain.Employee, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, age, position, address, created_at, updated_at FROM employee WHERE id = ?", id)
	var e domain.Employee
	if err := row.Scan(&e.ID, &e.Name, &e.Age, &e.Position, &e.Address, &e.CreatedAt, &e.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &e, nil
}

func (r *mysqlEmployeeRepo) Create(ctx context.Context, e *domain.Employee) error {
	res, err := r.db.ExecContext(ctx, "INSERT INTO employee (name, age, position, address, created_at, updated_at) VALUES (?,?,?,?,?,?)", e.Name, e.Age, e.Position, e.Address, e.CreatedAt, e.UpdatedAt)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	e.ID = int(id)
	return nil
}

func (r *mysqlEmployeeRepo) Update(ctx context.Context, e *domain.Employee) error {
	_, err := r.db.ExecContext(ctx, "UPDATE employee SET name=?, age=?, position=?, address=?, updated_at=? WHERE id=?", e.Name, e.Age, e.Position, e.Address, e.UpdatedAt, e.ID)
	return err
}

func (r *mysqlEmployeeRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM employee WHERE id=?", id)
	return err
}
