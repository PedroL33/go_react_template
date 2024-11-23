package store

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstHame    string `json:"first_name"`
	LastName     string `json:"last_name"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UserStore interface {
	CreateUser(email string, password string, firstName string, lastName string) error
}

type SQLUserStore struct {
	DB *sql.DB
}

func (s *SQLUserStore) CreateUser(
	email string,
	password string,
	firstName string,
	lastName string,
) error {

	var existingUser string
	err := s.DB.QueryRow("SELECT email FROM users WHERE email = $1", email).Scan(&existingUser)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	hashedPw := string(bytes)

	_, err = s.DB.Exec("INSERT INTO users (email, password_hash, first_name, last_name) VALUES ($1, $2, $3, $4)", email, hashedPw, firstName, lastName)

	return err
}
