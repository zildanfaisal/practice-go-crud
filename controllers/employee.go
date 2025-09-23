package controllers

import (
	"encoding/json"
	"net/http"
	"practice-go-crud/database"
	"practice-go-crud/models"
	"strconv"
	"time"
)

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT id, name, age, position, address, created_at, updated_at FROM employee")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.Position, &emp.Address, &emp.CreatedAt, &emp.UpdatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		employees = append(employees, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	row := database.DB.QueryRow("SELECT id, name, age, position, address, created_at, updated_at FROM employee WHERE id = ?", idStr)

	var emp models.Employee
	if err := row.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.Position, &emp.Address, &emp.CreatedAt, &emp.UpdatedAt); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()
	result, err := database.DB.Exec("INSERT INTO employee (name, age, position, address, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		emp.Name, emp.Age, emp.Position, emp.Address, now, now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	emp.ID = int(id)
	emp.CreatedAt = now
	emp.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	var emp models.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now := time.Now()
	_, err := database.DB.Exec("UPDATE employee SET name = ?, age = ?, position = ?, address = ?, updated_at = ? WHERE id = ?",
		emp.Name, emp.Age, emp.Position, emp.Address, now, idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := strconv.Atoi(idStr)
	emp.ID = id
	emp.UpdatedAt = now

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	_, err := database.DB.Exec("DELETE FROM employee WHERE id = ?", idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
