package main

import (
	"context"
	"errors"

	"github.com/Vadim-Makhnev/grpc/internal/data"
	"github.com/Vadim-Makhnev/grpc/internal/grpcutils"
	"github.com/Vadim-Makhnev/grpc/internal/validator"
	"github.com/Vadim-Makhnev/grpc/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	proto.UnimplementedUserServiceServer
	app *application
}

func (u *UserService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {

	user := &data.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		u.app.logger.Warn("validation failed", "errors", v.Errors)
		return nil, grpcutils.FailedValidation(v.Errors)
	}

	err := u.app.models.Users.CreateUser(user)
	if err != nil {
		return nil, grpcutils.Internal(u.app.logger, err, "")
	}

	resp := &proto.UserResponse{
		Id:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Age:     user.Age,
		Version: user.Version,
	}

	return resp, nil
}

func (u *UserService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.UserResponse, error) {
	id := req.Id
	if id <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user id must be greater than zero")
	}

	user, err := u.app.models.Users.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):

			return nil, grpcutils.NotFound(err.Error())
		default:
			return nil, grpcutils.Internal(u.app.logger, err, "")
		}
	}

	resp := &proto.UserResponse{
		Id:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Age:     user.Age,
		Version: user.Version,
	}

	return resp, nil
}

func (u *UserService) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	var input struct {
		data.Filters
	}

	v := validator.New()

	input.Page = int(u.app.getInt32(req.Page, 1))
	input.PageSize = int(u.app.getInt32(req.PageSize, 20))
	input.Sort = u.app.getString(req.Sort, "id")
	input.SortSafelist = []string{"id", "-id", "name", "email", "age"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		return nil, grpcutils.FailedValidation(v.Errors)
	}

	users, metadata, err := u.app.models.Users.GetAll(input.Filters)
	if err != nil {
		return nil, grpcutils.Internal(u.app.logger, err, "")
	}

	protoUsers := make([]*proto.UserResponse, len(users))
	for i, user := range users {
		protoUsers[i] = &proto.UserResponse{
			Id:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Age:     user.Age,
			Version: user.Version,
		}
	}

	protoMetadata := &proto.MetaData{
		TotalRecords: int32(metadata.TotalRecords),
		Page:         int32(metadata.CurrentPage),
		PageSize:     int32(metadata.PageSize),
	}

	resp := &proto.ListUsersResponse{
		Users:    protoUsers,
		Metadata: protoMetadata,
	}

	return resp, nil
}
