package domain

import "errors"

// Sentinel errors used across layers.
var (
	ErrNotFound          = errors.New("resource not found")
	ErrConflict          = errors.New("resource already exists")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInvalidInput      = errors.New("invalid input")
	ErrOTPExpired        = errors.New("OTP has expired")
	ErrOTPInvalid        = errors.New("OTP is invalid")
	ErrOTPAlreadyUsed    = errors.New("OTP has already been used")
	ErrUserNotVerified   = errors.New("user account is not verified")
	ErrRateLimitExceeded = errors.New("rate limit exceeded, please try again later")
	ErrInsufficientStock = errors.New("insufficient stock for one or more items")
	ErrEmptyCart         = errors.New("cart is empty")
)

// AppError wraps a sentinel error with an optional human-readable message.
type AppError struct {
	Err     error
	Message string
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(err error, msg string) *AppError {
	return &AppError{Err: err, Message: msg}
}
