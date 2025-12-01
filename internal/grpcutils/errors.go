package grpcutils

import (
	"log/slog"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrMessageUserNotFound           = "user not found"
	ErrMessageEditConflict           = "unable to update due to edit conflict, please try again"
	ErrMessageRateLimitExceeded      = "rate limit exceeded"
	ErrMessageInvalidCredentials     = "invalid authentication credentials"
	ErrMessageInvalidToken           = "invalid or missing authentication token"
	ErrMessageAuthenticationRequired = "you must be authenticated to access this resource"
	ErrMessageAccountInactive        = "your user account must be activated"
	ErrMessageNotPermitted           = "you don't have permission to access this resource"
	ErrMessageInternalProblem        = "the server encountered a problem"
	ErrMessageBadRequest             = "invalid request"
	ErrMessageInvalidRequest         = "invalid request"
)

func NotFound(msg string) error {
	if msg == "" {
		msg = ErrMessageUserNotFound
	}
	return status.Error(codes.NotFound, msg)
}

func FailedValidation(errors map[string]string) error {
	st := status.New(codes.InvalidArgument, ErrMessageInvalidRequest)

	for field, desc := range errors {
		violation := &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: desc,
		}
		stWithDetail, err := st.WithDetails(violation)
		if err != nil {
			return status.Error(codes.Internal, ErrMessageInternalProblem)
		}
		st = stWithDetail
	}

	return st.Err()
}

func Internal(logger *slog.Logger, err error, msg string) error {
	if msg == "" {
		msg = ErrMessageInternalProblem
	}

	if logger != nil {
		logger.Error("internal server error", "error", err)
	}

	return status.Error(codes.Internal, msg)
}
