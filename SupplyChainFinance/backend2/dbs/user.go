package dbs

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user record.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Address  string `json:"address"`
	Role     string `json:"role"`
}

var DB *sql.DB

// InitDB initializes database connection and required schema.
func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
	if err = EnsureSchema(); err != nil {
		log.Fatal(err)
	}
}

// CloseDB closes database connection.
func CloseDB() {
	if DB != nil {
		_ = DB.Close()
	}
}

// GetUser returns user info by username.
func GetUser(username string) (User, error) {
	var user User
	query := "SELECT id, username, password, address, role FROM users WHERE username = ?"
	err := DB.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Password, &user.Address, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, nil
		}
		return User{}, err
	}
	return user, nil
}

// ListUsers returns all users.
func ListUsers() ([]User, error) {
	rows, err := DB.Query("SELECT id, username, address, role FROM users ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err = rows.Scan(&u.ID, &u.Username, &u.Address, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// SaveUser inserts new user.
func SaveUser(username, password, address, role string) error {
	if strings.TrimSpace(role) == "" {
		role = "business"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := "INSERT INTO users (username, password, address, role) VALUES (?, ?, ?, ?)"
	_, err = DB.Exec(query, username, string(hashedPassword), address, role)
	if err != nil {
		return fmt.Errorf("save user failed: %w", err)
	}
	return nil
}
