package data

import (
	"testing"

	"github.com/Vadim-Makhnev/grpc/internal/validator"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateUser(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		wantErrs map[string]string
	}{
		{
			name: "valid user",
			user: &User{
				Name:  "Andrew",
				Email: "andrew@google.com",
				Age:   31,
			},
			wantErrs: nil,
		},
		{
			name: "empty name",
			user: &User{
				Name:  "",
				Email: "andrew@google.com",
				Age:   31,
			},
			wantErrs: map[string]string{
				"name": "must be provided",
			},
		},
		{
			name: "invalid age",
			user: &User{
				Name:  "Andrew",
				Email: "andrew@google.com",
				Age:   -5,
			},
			wantErrs: map[string]string{
				"age": "must be greater than 0",
			},
		},
		{
			name: "invalid email",
			user: &User{
				Name:  "Andrew",
				Email: "andrew.com",
				Age:   31,
			},
			wantErrs: map[string]string{
				"email": "must be a valid email address",
			},
		},
		{
			name: "empty email",
			user: &User{
				Name:  "Andrew",
				Email: "",
				Age:   31,
			},
			wantErrs: map[string]string{
				"email": "must be provided",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateUser(v, tt.user)

			if tt.wantErrs == nil {
				assert.True(t, v.Valid(), "expected no validation errors")
			} else {
				assert.False(t, v.Valid(), "expected validation errors")
				assert.Equal(t, tt.wantErrs, v.Errors)
			}
		})
	}
}
