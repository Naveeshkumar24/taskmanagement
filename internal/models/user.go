package models

import "net/http"

// User model for authentication and task assignment
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"` // omit in JSON response
	Role     string `json:"role"`     // e.g., admin, manager, user
}

// Task model representing the core task entity
type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"` // low, medium, high
	Status      string `json:"status"`   // todo, in-progress, done
	CreatedBy   int    `json:"created_by"`
	AssignedTo  int    `json:"assigned_to"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// TaskFilter struct for handling filter/search queries
type TaskFilter struct {
	Title      string `json:"title"`
	Status     string `json:"status"`
	Priority   string `json:"priority"`
	DueBefore  string `json:"due_before"`
	DueAfter   string `json:"due_after"`
	AssignedTo int    `json:"assigned_to"`
	CreatedBy  int    `json:"created_by"`
}

// Interfaces

type UserInterface interface {
	Register(user User) error
	Login(email, password string) (User, error)
	GetUserByID(id int) (User, error)
}

type TaskInterface interface {
	CreateTask(task Task) error
	GetTaskByID(id int) (Task, error)
	UpdateTask(task Task) error
	DeleteTask(id int) error
	ListTasks(r *http.Request) ([]Task, error)
	SearchAndFilterTasks(filter TaskFilter) ([]Task, error)
	GetUserDashboard(userID int) (map[string][]Task, error)
}
