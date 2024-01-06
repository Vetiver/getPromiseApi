package db

import (
	"context"
	"fmt"
	"os"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	pool *pgxpool.Pool
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string `json:"name"`
	Group 	 string `json:"group"`
	Email    string `json:"email"`
	Password string `json:"password"`
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