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

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo models.UserInterface) *UserHandler {
	return &UserHandler{
		userRepo: userRepo.(*repository.UserRepository),
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := utils.Decode(r, &user); err != nil {
		log.Printf("Register decode error: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.userRepo.Register(user); err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	utils.Encode(w, map[string]string{"message": "Registration successful"})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := utils.Decode(r, &creds); err != nil {
		log.Printf("Login decode error: %v", err)
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	// Validate user credentials and get the user object
	user, err := h.userRepo.Login(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	// 2) wrap into a single response object
	resp := map[string]interface{}{
		"user":  user,
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json")

	// Send the response with the user and token
	utils.Encode(w, resp)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user.Password = ""
	utils.Encode(w, user)
}
