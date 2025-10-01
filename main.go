package main

import (
	"log"
	"net/http"

	"practice-go-crud/database"
	"practice-go-crud/internal/handler"
	"practice-go-crud/internal/middleware"
	"practice-go-crud/internal/repository"
	"practice-go-crud/internal/service"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	employeeRepo := repository.NewMySQLEmployeeRepo(db)
	employeeSvc := service.NewEmployeeService(employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeSvc)

	userRepo := repository.NewMySQLUserRepo(db)
	tokenRepo := repository.NewMySQLTokenRepo(db)
	authSvc := service.NewAuthService(userRepo, tokenRepo)
	authHandler := handler.NewAuthHandler(authSvc)

	authMiddleware := middleware.NewAuthMiddleware(authSvc)

	http.HandleFunc("/auth/register", authHandler.Register)
	http.HandleFunc("/auth/login", authHandler.Login)
	http.HandleFunc("/auth/logout", authMiddleware.RequireAuth(authHandler.Logout))
	http.HandleFunc("/auth/me", authMiddleware.RequireAuth(authHandler.Me))

	http.HandleFunc("/auth/admin/register", authMiddleware.RequireRole("admin")(authHandler.AdminRegister))

	http.HandleFunc("/employees", authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			employeeHandler.List(w, r)
		case http.MethodPost:
			authMiddleware.RequireRole("admin")(employeeHandler.Create)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/employee", authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			employeeHandler.Get(w, r)
		case http.MethodPut, http.MethodPatch:
			authMiddleware.RequireRole("admin")(employeeHandler.Update)(w, r)
		case http.MethodDelete:
			authMiddleware.RequireRole("admin")(employeeHandler.Delete)(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	log.Println("Clean architecture server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
