package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"auth/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Close() error
	RegisterUser(models.User) error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	schema     = os.Getenv("DB_SCHEMA")
	dbInstance *service
)

// New creates a new instance of the database service.
func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

// RegisterUser inserts a new user into the database.
// It checks if the email already exists in the database.
// If the email does not exist, it inserts the user into the database.
// If the email already exists, it returns an error.
func (s *service) RegisterUser(user models.User) error {
	err := s.isLoginExists(user.Email)
	if err != nil {
		return err
	}

	query := `INSERT INTO users (login, passwordHash) VALUES ($1, $2)`
	_, err = s.db.Exec(query, user.Email, user.PasswordHash)
	return err
}

func (s *service) isLoginExists(email string) error {
	query := `SELECT login FROM users WHERE login = $1`
	rows, err := s.db.Query(query, email)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return fmt.Errorf("email already exists")
	}

	return nil
}
