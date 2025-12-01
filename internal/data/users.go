package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Vadim-Makhnev/grpc/internal/validator"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int32     `json:"age"`
	CreatedAt time.Time `json:"-"`
	Version   int32     `json:"version"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) CreateUser(user *User) error {
	query := `
		INSERT INTO users (name, email, age)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version
		`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	args := []any{user.Name, user.Email, user.Age}

	return u.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
}

func (u *UserModel) GetUser(id int64) (*User, error) {
	query := `
		SELECT id, name, email, age, version
		FROM users
		WHERE id = $1
		`
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil

}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")

	v.Check(user.Email != "", "email", "must be provided")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be a valid email address")
}
