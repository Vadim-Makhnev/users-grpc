package mocks

import "github.com/Vadim-Makhnev/grpc/internal/data"

type UserStorageMock struct{}

func NewUserStorageMock() UserStorageMock {
	return UserStorageMock{}
}

func (s UserStorageMock) CreateUser(user *data.User) error {
	user.ID = 1
	user.Version = 1

	return nil
}

func (s UserStorageMock) GetUser(id int64) (*data.User, error) {
	if id == 1 {
		return &data.User{
			ID:      1,
			Name:    "Andrew",
			Email:   "andrew@google.com",
			Age:     31,
			Version: 1,
		}, nil
	}
	if id == 3 {
		return &data.User{
			ID:      3,
			Name:    "Andrew",
			Email:   "andrew@google.com",
			Age:     31,
			Version: 1,
		}, nil
	} else {
		return nil, data.ErrRecordNotFound
	}
}

func (s UserStorageMock) GetAll(filters data.Filters) ([]*data.User, data.MetaData, error) {
	return []*data.User{
		{
			ID:      1,
			Name:    "Andrew",
			Email:   "andrew@google.com",
			Age:     31,
			Version: 1,
		},
		{
			ID:      2,
			Name:    "Andrew",
			Email:   "andrew@google.com",
			Age:     31,
			Version: 1,
		},
	}, data.MetaData{}, nil
}

func (s UserStorageMock) DeleteUserById(id int64) (*data.User, error) {
	if id == 1 {
		return &data.User{
			ID:      1,
			Name:    "Andrew",
			Email:   "andrew@google.com",
			Age:     31,
			Version: 1,
		}, nil
	} else {
		return nil, data.ErrRecordNotFound
	}
}

func (s UserStorageMock) UpdateUser(user *data.User) error {
	if user.ID == 1 {
		user.Version = user.Version + 1
		return nil
	} else {
		return data.ErrEditConflict
	}
}
