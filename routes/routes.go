package routes

import (
	"net/http"
	"practice-go-crud/controllers"
)

func SetupRoutes() {
	http.HandleFunc("/employees", controllers.GetEmployees)
	http.HandleFunc("/employee", controllers.GetEmployeeByID)
	http.HandleFunc("/employee/create", controllers.CreateEmployee)
	http.HandleFunc("/employee/update", controllers.UpdateEmployee)
	http.HandleFunc("/employee/delete", controllers.DeleteEmployee)
}
