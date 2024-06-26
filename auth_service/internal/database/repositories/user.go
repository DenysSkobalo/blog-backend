package repositories

import (
	"auth_service/internal/api/models"
	"database/sql"
	"log"
)

type UserRepository interface {
	CreateUser(u *models.User) error
	GetUserByUsername(username string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

type Repositories struct {
	UserRepo UserRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		UserRepo: NewUserRepository(db),
	}
}

type SQLUserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *SQLUserRepository {
	return &SQLUserRepository{DB: db}
}

func (repo *SQLUserRepository) CreateUser(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.HashPassword(); err != nil {
		log.Printf("Unable to hash password: %v", err)
		return err
	}

	const query = `INSERT INTO users (username, email, password_hash, created_at, first_name, last_name) VALUES ($1, $2, $3, NOW(), $4, $5)`
	_, err := repo.DB.Exec(query, u.Username, u.Email, u.Password, u.FirstName, u.LastName)
	if err != nil {
		log.Printf("Unable to create user: %v", err)
		return err
	}
	return nil
}

func (repo *SQLUserRepository) GetUserByUsername(username string) (*models.User, error) {
	const query = `SELECT * FROM users WHERE username = $1`
	var user models.User
	err := repo.DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *SQLUserRepository) GetUserByEmail(email string) (*models.User, error) {
	const query = `SELECT * FROM users WHERE email = $1`
	var user models.User
	err := repo.DB.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
