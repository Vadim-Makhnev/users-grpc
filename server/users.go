package main

import (
	"context"
	"errors"

	"github.com/Vadim-Makhnev/grpc/internal/data"
	"github.com/Vadim-Makhnev/grpc/internal/grpcutils"
	"github.com/Vadim-Makhnev/grpc/internal/validator"
	"github.com/Vadim-Makhnev/grpc/proto"
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

	user, err := u.app.models.Users.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, grpcutils.NotFound("")
		case errors.Is(err, data.ErrInvalidArgument):
			return nil, grpcutils.InvalidArgument(u.app.logger, err, "")
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

func (u *UserService) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.UserResponse, error) {
	id := req.Id

	user, err := u.app.models.Users.DeleteUserById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, grpcutils.NotFound("")
		case errors.Is(err, data.ErrInvalidArgument):
			return nil, grpcutils.InvalidArgument(u.app.logger, err, "")
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

func (u *UserService) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UserResponse, error) {
	id := req.Id

	user, err := u.app.models.Users.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			return nil, grpcutils.NotFound("")
		default:
			return nil, grpcutils.Internal(u.app.logger, err, "")
		}
	}

	if req.Name != nil {
		user.Name = req.Name.Value
	}

	if req.Email != nil {
		user.Email = req.Email.Value
	}

	if req.Age != nil {
		user.Age = req.Age.Value
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		return nil, grpcutils.FailedValidation(v.Errors)
	}

	err = u.app.models.Users.UpdateUser(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			return nil, grpcutils.EditConflict(u.app.logger, err, "")
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
