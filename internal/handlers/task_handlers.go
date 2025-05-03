package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/naveeshkumar24/internal/models"
	"github.com/naveeshkumar24/pkg/utils"
	"github.com/naveeshkumar24/repository"
)

type TaskHandler struct {
	taskRepo *repository.TaskRepository
}

func NewTaskHandler(taskRepo models.TaskInterface) *TaskHandler {
	return &TaskHandler{
		taskRepo: taskRepo.(*repository.TaskRepository),
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := utils.Decode(r, &task)
	if err != nil {
		log.Printf("Failed to decode task data: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}

	err = h.taskRepo.CreateTask(task)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to create task"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	utils.Encode(w, map[string]string{"message": "Task created successfully"})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid task ID"})
		return
	}

	task, err := h.taskRepo.GetTaskByID(id)
	if err != nil {
		log.Printf("Task not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		utils.Encode(w, map[string]string{"message": "Task not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := utils.Decode(r, &task)
	if err != nil {
		log.Printf("Failed to decode task update: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid request body"})
		return
	}

	err = h.taskRepo.UpdateTask(task)
	if err != nil {
		log.Printf("Failed to update task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to update task"})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "Task updated successfully"})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID for deletion: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid task ID"})
		return
	}

	err = h.taskRepo.DeleteTask(id)
	if err != nil {
		log.Printf("Failed to delete task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to delete task"})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, map[string]string{"message": "Task deleted successfully"})
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.ListTasks(r)
	if err != nil {
		log.Printf("Failed to list tasks: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to list tasks"})
		return
	}

	if len(tasks) == 0 {
		log.Printf("No tasks found")
		w.WriteHeader(http.StatusNotFound)
		utils.Encode(w, map[string]string{"message": "No tasks found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, tasks)
}
func (h *TaskHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userID"]
	if userIDStr == "" {
		log.Println("Missing userID in path parameters")
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Missing userID in path parameter"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid userID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		utils.Encode(w, map[string]string{"message": "Invalid userID"})
		return
	}

	dashboard, err := h.taskRepo.GetUserDashboard(userID)
	if err != nil {
		log.Printf("Failed to get dashboard data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		utils.Encode(w, map[string]string{"message": "Failed to get dashboard data"})
		return
	}

	w.WriteHeader(http.StatusOK)
	utils.Encode(w, dashboard)
}
