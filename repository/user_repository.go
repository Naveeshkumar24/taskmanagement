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

func (u *UserRepository) Login(email, password string) (models.User, error) {
	query := database.NewQuery(u.db)
	user, err := query.GetUserByEmail(email)
	if err != nil {
		log.Printf("Repository: Failed to get user by email: %v", err)
		return models.User{}, err
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Printf("Repository: Incorrect password for user %s", email)
		return models.User{}, sql.ErrNoRows
	}

	return user, nil
}

func (u *UserRepository) GetUserByID(id int) (models.User, error) {
	query := database.NewQuery(u.db)
	user, err := query.GetUserByID(id)
	if err != nil {
		log.Printf("Repository: Failed to get user by ID: %v", err)
		return models.User{}, err
	}
	return user, nil
}
