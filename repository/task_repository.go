package repository

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/naveeshkumar24/internal/models"
	"github.com/naveeshkumar24/pkg/database"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (t *TaskRepository) CreateTask(task models.Task) error {
	query := database.NewQuery(t.db)
	err := query.CreateTask(task)
	if err != nil {
		log.Printf("Repository: Failed to create task: %v", err)
		return err
	}
	return nil
}

func (t *TaskRepository) GetTaskByID(id int) (models.Task, error) {
	query := database.NewQuery(t.db)
	task, err := query.GetTaskByID(id)
	if err != nil {
		log.Printf("Repository: Failed to fetch task by ID: %v", err)
		return models.Task{}, err
	}
	return task, nil
}

func (t *TaskRepository) UpdateTask(task models.Task) error {
	query := database.NewQuery(t.db)
	err := query.UpdateTask(task)
	if err != nil {
		log.Printf("Repository: Failed to update task: %v", err)
		return err
	}
	return nil
}

func (t *TaskRepository) DeleteTask(id int) error {
	query := database.NewQuery(t.db)
	err := query.DeleteTask(id)
	if err != nil {
		log.Printf("Repository: Failed to delete task: %v", err)
		return err
	}
	return nil
}

func (t *TaskRepository) ListTasks(r *http.Request) ([]models.Task, error) {
	query := database.NewQuery(t.db)
	tasks, err := query.ListTasks()
	if err != nil {
		log.Printf("Repository: Failed to list tasks: %v", err)
		return nil, err
	}
	return tasks, nil
}

func (t *TaskRepository) SearchAndFilterTasks(filter models.TaskFilter) ([]models.Task, error) {
	query := database.NewQuery(t.db)
	tasks, err := query.SearchAndFilterTasks(filter)
	if err != nil {
		log.Printf("Repository: Failed to search/filter tasks: %v", err)
		return nil, err
	}
	return tasks, nil
}

func (t *TaskRepository) GetUserDashboard(userID int) (map[string][]models.Task, error) {
	query := database.NewQuery(t.db)
	dashboardData, err := query.GetUserDashboard(userID)
	if err != nil {
		log.Printf("Repository: Failed to get user dashboard: %v", err)
		return nil, err
	}
	return dashboardData, nil
}
