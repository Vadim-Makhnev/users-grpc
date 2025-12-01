package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Vadim-Makhnev/grpc/internal/validator"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Age       int32
	CreatedAt time.Time
	Version   int32
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
		SELECT id, name, email, age, created_at, version
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
		&user.CreatedAt,
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

func (u *UserModel) GetAll(filters Filters) ([]*User, MetaData, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, name, email, age, created_at, version
		FROM users
		ORDER BY %s %s, id ASC
		LIMIT $1 OFFSET $2`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{filters.limit(), filters.offset()}

	rows, err := u.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, MetaData{}, err
	}

	defer rows.Close()

	totalRecords := 0

	users := []*User{}

	for rows.Next() {
		var user User

		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
			&user.Version,
		)
		if err != nil {
			return nil, MetaData{}, err
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return users, metadata, err
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(user.Email != "", "email", "must be provided")
	v.Check(user.Age > 0, "age", "must be greater than 0")
	v.Check(validator.Matches(user.Email, validator.EmailRX), "email", "must be a valid email address")
}
