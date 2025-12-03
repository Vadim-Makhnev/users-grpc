package main

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/Vadim-Makhnev/grpc/internal/data"
	"github.com/Vadim-Makhnev/grpc/internal/data/mocks"
	"github.com/Vadim-Makhnev/grpc/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUserService_CreateUser_Success(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}

	UserService := &UserService{app: app}

	ctx := context.Background()

	user, err := UserService.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "Andrew",
		Email: "andrew@google.com",
		Age:   31,
	})
	assert.NoError(t, err, "no error expected")
	assert.Equal(t, "Andrew", user.Name)
	assert.Equal(t, "andrew@google.com", user.Email)
	assert.Equal(t, int32(31), user.Age)
	assert.NotEmpty(t, user.Id)
	assert.NotEmpty(t, user.Version)
}

func TestUserService_CreateUser_InvalidEmail(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	_, err := service.CreateUser(context.Background(), &proto.CreateUserRequest{
		Name:  "Andrew",
		Email: "invalid-email",
		Age:   31,
	})

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUserService_GetUser_Success(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	ctx := context.Background()

	user, err := service.GetUser(ctx, &proto.GetUserRequest{
		Id: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.Id)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Age)
	assert.NotEmpty(t, user.Version)
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	ctx := context.Background()

	_, err := service.GetUser(ctx, &proto.GetUserRequest{
		Id: 2,
	})

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestUserService_DeleteUser_Success(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	ctx := context.Background()

	user, err := service.DeleteUser(ctx, &proto.DeleteUserRequest{
		Id: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.Id)
	assert.NotEmpty(t, user.Name)
	assert.NotEmpty(t, user.Email)
	assert.NotEmpty(t, user.Age)
	assert.NotEmpty(t, user.Version)
}

func TestUserService_DeleteUser_NotFound(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	ctx := context.Background()

	_, err := service.DeleteUser(ctx, &proto.DeleteUserRequest{
		Id: 2,
	})

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestUserService_UpdateUser_Success(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	ctx := context.Background()

	user, err := service.UpdateUser(ctx, &proto.UpdateUserRequest{
		Id:    1,
		Name:  wrapperspb.String("John"),
		Email: wrapperspb.String("john@gmail.com"),
		Age:   wrapperspb.Int32(21),
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.Id)
	assert.Equal(t, "John", user.Name)
	assert.Equal(t, "john@gmail.com", user.Email)
	assert.Equal(t, int32(21), user.Age)
	assert.Equal(t, int32(2), user.Version)
}

func TestUserService_UpdateUser_EditConflict(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	ctx := context.Background()

	_, err := service.UpdateUser(ctx, &proto.UpdateUserRequest{
		Id:    3,
		Name:  wrapperspb.String("John"),
		Email: wrapperspb.String("john@gmail.com"),
		Age:   wrapperspb.Int32(21),
	})

	assert.Error(t, err)
	st, _ := status.FromError(err)
	assert.Equal(t, codes.Aborted, st.Code())
}

func TestUserService_GetAll_Success(t *testing.T) {
	models := data.Models{
		Users: mocks.NewUserStorageMock(),
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))

	app := &application{
		logger: logger,
		models: models,
	}
	service := &UserService{app: app}

	ctx := context.Background()

	users, err := service.ListUsers(ctx, &proto.ListUsersRequest{})

	assert.NoError(t, err)
	assert.NotEmpty(t, users)
}
