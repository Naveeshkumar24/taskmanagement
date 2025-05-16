package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/naveeshkumar24/internal/models"
)

type Query struct {
	db   *sql.DB
	Time *time.Location
}

func NewQuery(db *sql.DB) *Query {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatalf("Failed to load time zone: %v", err)
	}
	return &Query{
		db:   db,
		Time: loc,
	}
}

func (q *Query) CreateTaskTables() error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(100) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password TEXT NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user'
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			due_date DATE,
			priority VARCHAR(50),
			status VARCHAR(50),
			created_by INT REFERENCES users(id),
			assigned_to INT REFERENCES users(id),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := tx.Exec(query); err != nil {
			log.Printf("Failed to execute query: %s", query)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	log.Println("User and Task tables created successfully.")
	return nil
}

// ======================== User Functions ========================

func (q *Query) RegisterUser(user models.User) error {
	// Check if the username already exists

	// Proceed to insert the new user if no duplicates were found
	_, err := q.db.Exec(`
        INSERT INTO users (username, email, password, role) 
        VALUES ($1, $2, $3, $4)
    `, user.Username, user.Email, user.Password, user.Role)

	if err != nil {
		log.Printf("Failed to register user: %v", err)
		return fmt.Errorf("error registering user: %w", err)
	}

	return nil
}

func (q *Query) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	err := q.db.QueryRow(`
		SELECT id, username, email, password, role FROM users WHERE email = $1
	`, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)

	return user, err
}

func (q *Query) GetUserByID(id int) (models.User, error) {
	var user models.User

	err := q.db.QueryRow(`
		SELECT id, username, email, password, role FROM users WHERE id = $1
	`, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)

	return user, err
}

// ======================== Task Functions ========================

func (q *Query) CreateTask(task models.Task) error {
	// Check if the user exists
	var userCount int
	err := q.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", task.CreatedBy).Scan(&userCount)
	if err != nil {
		log.Printf("Failed to check user existence: %v", err)
		return err
	}

	if userCount == 0 {
		return fmt.Errorf("user with ID %d does not exist", task.CreatedBy)
	}

	// Proceed with task insertion
	_, err = q.db.Exec(`
        INSERT INTO tasks (title, description, due_date, priority, status, created_by, assigned_to)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, task.Title, task.Description, task.DueDate, task.Priority, task.Status, task.CreatedBy, task.AssignedTo)

	if err != nil {
		log.Printf("Failed to create task: %v", err)
		return err
	}
	log.Printf("Task %s created successfully.", task.Title)
	return nil
}

func (q *Query) GetTaskByID(id int) (models.Task, error) {
	var task models.Task

	err := q.db.QueryRow(`
		SELECT id, title, description, due_date, priority, status, created_by, assigned_to, created_at, updated_at
		FROM tasks WHERE id = $1
	`, id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Priority, &task.Status, &task.CreatedBy, &task.AssignedTo, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		log.Printf("Failed to fetch task by ID: %v", err)
		return task, err
	}

	return task, nil
}

func (q *Query) UpdateTask(task models.Task) error {
	_, err := q.db.Exec(`
		UPDATE tasks SET
			title = $1, description = $2, due_date = $3, priority = $4, status = $5,
			assigned_to = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
	`, task.Title, task.Description, task.DueDate, task.Priority, task.Status, task.AssignedTo, task.ID)

	if err != nil {
		log.Printf("Failed to update task ID %d: %v", task.ID, err)
		return err
	}

	log.Printf("Task ID %d updated successfully.", task.ID)
	return nil
}

func (q *Query) DeleteTask(id int) error {
	_, err := q.db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Failed to delete task ID %d: %v", id, err)
		return err
	}
	log.Printf("Task ID %d deleted successfully.", id)
	return nil
}

func (q *Query) ListTasks() ([]models.Task, error) {
	var tasks []models.Task

	rows, err := q.db.Query(`
		SELECT id, title, description, due_date, priority, status, created_by, assigned_to, created_at, updated_at
		FROM tasks
	`)
	if err != nil {
		log.Printf("Failed to list tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Priority, &task.Status,
			&task.CreatedBy, &task.AssignedTo, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			log.Printf("Failed to scan task row: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
func (q *Query) GetUserDashboard(userID int) (map[string][]models.Task, error) {
	dashboard := make(map[string][]models.Task)

	rows, err := q.db.Query(`
		SELECT id, title, description, due_date, priority, status, created_by, assigned_to, created_at, updated_at
		FROM tasks
		WHERE created_by = $1 OR assigned_to = $1
		ORDER BY status
	`, userID)
	if err != nil {
		log.Printf("Failed to get dashboard data for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Priority,
			&task.Status, &task.CreatedBy, &task.AssignedTo, &task.CreatedAt, &task.UpdatedAt); err != nil {
			log.Printf("Failed to scan dashboard task row: %v", err)
			return nil, err
		}

		status := task.Status
		dashboard[status] = append(dashboard[status], task)
	}

	log.Printf("Dashboard data fetched successfully for user ID %d", userID)
	return dashboard, nil
}
func (q *Query) SearchAndFilterTasks(filter models.TaskFilter) ([]models.Task, error) {
	var tasks []models.Task

	query := `SELECT id, title, description, due_date, priority, status, created_by, assigned_to, created_at, updated_at
	          FROM tasks WHERE 1=1`
	args := []interface{}{}
	argID := 1

	if filter.Status != "" {
		query += " AND status = $" + strconv.Itoa(argID)
		args = append(args, filter.Status)
		argID++
	}
	if filter.Priority != "" {
		query += " AND priority = $" + strconv.Itoa(argID)
		args = append(args, filter.Priority)
		argID++
	}
	if filter.AssignedTo != 0 {
		query += " AND assigned_to = $" + strconv.Itoa(argID)
		args = append(args, filter.AssignedTo)
		argID++
	}

	rows, err := q.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Priority,
			&task.Status, &task.CreatedBy, &task.AssignedTo, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
