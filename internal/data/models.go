package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrEditConflict    = errors.New("edit conflict")
	ErrInvalidArgument = errors.New("invalid argument")
)

type UserStorage interface {
	CreateUser(user *User) error
	GetUser(id int64) (*User, error)
	GetAll(filters Filters) ([]*User, MetaData, error)
	DeleteUserById(id int64) (*User, error)
	UpdateUser(user *User) error
}

type Models struct {
	Users UserStorage
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users: UserModel{
			DB: db,
		},
	}
}
