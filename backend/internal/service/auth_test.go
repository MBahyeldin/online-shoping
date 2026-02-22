package service_test

import (
	"testing"

	"github.com/online-cake-shop/backend/internal/service"
)

func TestGenerateOTP(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{"6 digits", 6},
		{"8 digits", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := service.GenerateOTP(tt.length)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(otp) != tt.length {
				t.Errorf("expected length %d, got %d (value: %s)", tt.length, len(otp), otp)
			}
			for _, c := range otp {
				if c < '0' || c > '9' {
					t.Errorf("non-digit character in OTP: %c", c)
				}
			}
		})
	}
}

func TestValidateRegisterInput(t *testing.T) {
	tests := []struct {
		name    string
		input   service.RegisterInput
		wantErr bool
	}{
		{
			name: "valid input",
			input: service.RegisterInput{
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "+12025551234",
				Email:       "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing first name",
			input: service.RegisterInput{
				FirstName:   "",
				LastName:    "Doe",
				PhoneNumber: "+12025551234",
				Email:       "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			input: service.RegisterInput{
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "+12025551234",
				Email:       "not-an-email",
			},
			wantErr: true,
		},
		{
			name: "invalid phone – too short",
			input: service.RegisterInput{
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "123",
				Email:       "john@example.com",
			},
			wantErr: true,
		},
		{
			name: "empty email",
			input: service.RegisterInput{
				FirstName:   "John",
				LastName:    "Doe",
				PhoneNumber: "+12025551234",
				Email:       "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateRegisterInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr=%v, got err=%v", tt.wantErr, err)
			}
		})
	}
}
