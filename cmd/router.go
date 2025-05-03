package main

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/naveeshkumar24/internal/handlers"
	"github.com/naveeshkumar24/internal/middleware"
	"github.com/naveeshkumar24/repository"
)

// registerTaskRouter sets up and configures the HTTP routes for task management operations.
// It initializes a new router with CORS middleware and defines endpoints for creating,
// fetching, updating, deleting, listing, and getting the dashboard view of tasks.
// The function takes a database connection as input and returns a configured router.
func registerTaskRouter(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	router.Use(middleware.CorsMiddleware)

	taskRepo := repository.NewTaskRepository(db)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Task routes
	router.HandleFunc("/task/create", taskHandler.CreateTask).Methods("POST")
	router.HandleFunc("/task/get/{id}", taskHandler.GetTask).Methods("GET")
	router.HandleFunc("/task/update", taskHandler.UpdateTask).Methods("POST")
	router.HandleFunc("/task/delete/{id}", taskHandler.DeleteTask).Methods("POST")
	router.HandleFunc("/task/list", taskHandler.ListTasks).Methods("GET")
	router.HandleFunc("/task/dashboard/{userID}", taskHandler.GetDashboard).Methods("GET")

	// User routes
	router.HandleFunc("/user/register", userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/user/login", userHandler.LoginUser).Methods("POST")
	router.HandleFunc("/user/get/{id}", userHandler.GetUserByID).Methods("GET")

	return router
}
