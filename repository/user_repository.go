package repository

import (
	"database/sql"
	"log"

	"github.com/naveeshkumar24/internal/models"
	"github.com/naveeshkumar24/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Register - Hashes the password and saves the user to the database
func (u *UserRepository) Register(user models.User) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return err
	}
	user.Password = string(hashedPassword)

	query := database.NewQuery(u.db)
	err = query.RegisterUser(user)
	if err != nil {
		log.Printf("Repository: Failed to register user: %v", err)
		return err
	}
	return nil
}

// Login - Verifies user credentials and returns user data
func (u *UserRepository) Login(email, password string) (models.User, error) {
	query := database.NewQuery(u.db)
	user, err := query.GetUserByEmail(email)
	if err != nil {
		// Check if the error is sql.ErrNoRows (when user is not found)
		if err == sql.ErrNoRows {
			log.Printf("Repository: User with email %s not found", email)
		} else {
			log.Printf("Repository: Failed to get user by email: %v", err)
		}
		return models.User{}, err
	}

	// Compare the hashed password with the entered password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Provide specific log message for password mismatch
		log.Printf("Repository: Incorrect password for user %s", email)
		return models.User{}, sql.ErrNoRows
	}

	return user, nil
}

// GetUserByID - Retrieves user details by ID
func (u *UserRepository) GetUserByID(id int) (models.User, error) {
	query := database.NewQuery(u.db)
	user, err := query.GetUserByID(id)
	if err != nil {
		log.Printf("Repository: Failed to get user by ID: %v", err)
		return models.User{}, err
	}
	return user, nil
}
