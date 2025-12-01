package main

import (
	"context"
	"errors"

	"github.com/Vadim-Makhnev/grpc/internal/data"
	"github.com/Vadim-Makhnev/grpc/internal/validator"
	"github.com/Vadim-Makhnev/grpc/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
		badRequest := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{},
		}

		for field, message := range v.Errors {
			badRequest.FieldViolations = append(badRequest.FieldViolations, &errdetails.BadRequest_FieldViolation{
				Field:       field,
				Description: message,
			})
		}

		statusProto := status.New(codes.InvalidArgument, "invalid request")
		statusWithDetails, err := statusProto.WithDetails(badRequest)
		if err != nil {
			return nil, status.Error(codes.Internal, "failed to construct error details")
		}

		u.app.logger.Warn("field error", "errors", v.Errors)
		return nil, statusWithDetails.Err()
	}

	err := u.app.models.Users.CreateUser(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
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
