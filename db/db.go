package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type DB struct {
	pool *pgxpool.Pool
}

type User struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"name"     binding:"required"`
	Group       string    `json:"group"    binding:"required"`
	Email       string    `json:"email"    binding:"required"`
	Password    string    `json:"password" binding:"required,min=8"`
	ConfirmCode int
}

type UserLoginData struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Token struct {
	TokenString string `json:"accessToken"`
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func DbStart(baseUrl string) *pgxpool.Pool {
	urlExample := baseUrl
	dbpool, err := pgxpool.New(context.Background(), string(urlExample))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v", err)
		os.Exit(1)
	}
	return dbpool
}

func (db DB) RegisterUser(userData User) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	userData.ID = uuid.New()
	password, hashErr := hashPassword(userData.Password)
	if hashErr != nil {
		return nil, fmt.Errorf("unable to hashPass: %v", hashErr)
	}

	err = conn.QueryRow(context.Background(),
		"INSERT INTO users(id, name, \"group\", email, password) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		userData.ID, userData.Username, userData.Group, userData.Email, password).Scan(&userData.ID)
	if err != nil {
		return nil, fmt.Errorf("unable to INSERT: %v", err)
	}

	return &userData, nil
}

func (db DB) GetUserByEmail(email string) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	var user User
	err = conn.QueryRow(context.Background(), "SELECT id, name, \"group\", email, password FROM users WHERE email = $1", email).
		Scan(&user.ID, &user.Username, &user.Group, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve user: %v", err)
	}

	return &user, err
}

func (db DB) userExists(userID uuid.UUID) (bool, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return false, fmt.Errorf("unable to acquire a database connection: %v", err)
	}
	defer conn.Release()

	var exists bool
	err = conn.QueryRow(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", userID).
		Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking user existence: %v", err)
	}

	return exists, nil
}

func (db DB) GetAllUsers(userID uuid.UUID) ([]User, error) {
    exists, err := db.userExists(userID)
    if err != nil {
        return nil, err
    }

    if exists == false {
        return nil, fmt.Errorf("user with ID %s does not exist", userID.String())
    }

    conn, err := db.pool.Acquire(context.Background())
    if err != nil {
        return nil, fmt.Errorf("unable to acquire a database connection: %v", err)
    }
    defer conn.Release()

    rows, err := conn.Query(context.Background(),
        "SELECT \"name\", \"group\", email FROM users")
    if err != nil {
        return nil, fmt.Errorf("unable to retrieve data from database: %v", err)
    }
    defer rows.Close()

    var data []User
    for rows.Next() {
        var d User
        err = rows.Scan(&d.Username, &d.Group, &d.Email )
        if err != nil {
            return nil, fmt.Errorf("unable to scan row: %v", err)
        }
        data = append(data, d)
    }
    return data, err
}