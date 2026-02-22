package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/online-cake-shop/backend/internal/config"
	"github.com/online-cake-shop/backend/internal/domain"
	"github.com/online-cake-shop/backend/internal/email"
	"github.com/online-cake-shop/backend/internal/repository/db"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{7,14}$`)
)

const (
	otpLength      = 6
	maxOTPAttempts = 5
	maxOTPPerHour  = 3
)

type AuthService struct {
	q         *db.Queries
	emailSvc  email.Sender
	jwtConfig config.JWTConfig
}

func NewAuthService(q *db.Queries, emailSvc email.Sender, jwtConfig config.JWTConfig) *AuthService {
	return &AuthService{q: q, emailSvc: emailSvc, jwtConfig: jwtConfig}
}

// ─── DTOs ────────────────────────────────────────────────────────────────────

type RegisterInput struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

type VerifyOTPInput struct {
	Email string
	OTP   string
}

type AuthResult struct {
	Token string
	User  db.User
}

// ─── Register ────────────────────────────────────────────────────────────────

func (s *AuthService) Register(ctx context.Context, in RegisterInput) error {
	if err := validateRegisterInput(in); err != nil {
		return err
	}

	// Check for duplicate email
	existing, err := s.q.GetUserByEmail(ctx, strings.ToLower(in.Email))
	if err == nil && existing.IsVerified {
		return &domain.AppError{Err: domain.ErrConflict, Message: "an account with this email already exists"}
	}

	var user db.User
	if err == nil {
		// Unverified user exists — resend OTP
		user = existing
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("get user by email: %w", err)
	} else {
		// Check phone uniqueness
		if _, err := s.q.GetUserByPhone(ctx, in.PhoneNumber); err == nil {
			return &domain.AppError{Err: domain.ErrConflict, Message: "an account with this phone number already exists"}
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("get user by phone: %w", err)
		}

		// Create new user
		user, err = s.q.CreateUser(ctx, db.CreateUserParams{
			FirstName:    in.FirstName,
			LastName:     in.LastName,
			PhoneNumber:  in.PhoneNumber,
			EmailAddress: strings.ToLower(in.Email),
		})
		if err != nil {
			return fmt.Errorf("create user: %w", err)
		}
	}

	return s.sendOTP(ctx, user)
}

// ─── Verify OTP ──────────────────────────────────────────────────────────────

func (s *AuthService) VerifyOTP(ctx context.Context, in VerifyOTPInput) (*AuthResult, error) {
	user, err := s.q.GetUserByEmail(ctx, strings.ToLower(in.Email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.AppError{Err: domain.ErrNotFound, Message: "no account found with this email"}
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	otp, err := s.q.GetLatestOTPByUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &domain.AppError{Err: domain.ErrOTPInvalid, Message: "no OTP found, please request a new one"}
		}
		return nil, fmt.Errorf("get otp: %w", err)
	}

	if otp.IsUsed {
		return nil, domain.ErrOTPAlreadyUsed
	}
	if time.Now().After(otp.ExpiresAt) {
		return nil, domain.ErrOTPExpired
	}
	if otp.AttemptCount >= maxOTPAttempts {
		return nil, &domain.AppError{Err: domain.ErrOTPInvalid, Message: "too many failed attempts, please request a new OTP"}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(otp.OtpHash), []byte(in.OTP)); err != nil {
		// Increment attempt counter but don't reveal the mismatch in error detail
		if _, incErr := s.q.IncrementOTPAttempts(ctx, otp.ID); incErr != nil {
			return nil, fmt.Errorf("increment otp attempts: %w", incErr)
		}
		return nil, domain.ErrOTPInvalid
	}

	// Mark OTP used
	if err := s.q.MarkOTPUsed(ctx, otp.ID); err != nil {
		return nil, fmt.Errorf("mark otp used: %w", err)
	}

	// Mark user verified
	verifiedUser, err := s.q.MarkUserVerified(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("mark user verified: %w", err)
	}

	token, err := s.generateJWT(verifiedUser.ID)
	if err != nil {
		return nil, fmt.Errorf("generate jwt: %w", err)
	}

	return &AuthResult{Token: token, User: verifiedUser}, nil
}

// ─── Resend OTP ──────────────────────────────────────────────────────────────

func (s *AuthService) ResendOTP(ctx context.Context, email string) error {
	user, err := s.q.GetUserByEmail(ctx, strings.ToLower(email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &domain.AppError{Err: domain.ErrNotFound, Message: "no account found with this email"}
		}
		return fmt.Errorf("get user: %w", err)
	}
	if user.IsVerified {
		return &domain.AppError{Err: domain.ErrConflict, Message: "account is already verified"}
	}

	return s.sendOTP(ctx, user)
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

func (s *AuthService) sendOTP(ctx context.Context, user db.User) error {
	// Rate limit: max 3 OTPs per hour per user
	count, err := s.q.CountRecentOTPsByUserID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("count otps: %w", err)
	}
	if count >= maxOTPPerHour {
		return domain.ErrRateLimitExceeded
	}

	// Invalidate previous OTPs
	if err := s.q.InvalidateUserOTPs(ctx, user.ID); err != nil {
		return fmt.Errorf("invalidate otps: %w", err)
	}

	// Generate OTP
	rawOTP, err := generateOTP(otpLength)
	if err != nil {
		return fmt.Errorf("generate otp: %w", err)
	}

	// Hash OTP
	hash, err := bcrypt.GenerateFromPassword([]byte(rawOTP), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash otp: %w", err)
	}

	// Store OTP
	if _, err := s.q.CreateOTP(ctx, db.CreateOTPParams{
		UserID:    user.ID,
		OtpHash:   string(hash),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}); err != nil {
		return fmt.Errorf("store otp: %w", err)
	}

	// Send email
	return s.emailSvc.SendOTP(user.EmailAddress, user.FirstName, rawOTP)
}

func (s *AuthService) generateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(s.jwtConfig.AccessTokenTTL).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtConfig.Secret))
}

func generateOTP(length int) (string, error) {
	digits := make([]byte, length)
	for i := range digits {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		digits[i] = byte('0' + n.Int64())
	}
	return string(digits), nil
}

func validateRegisterInput(in RegisterInput) error {
	in.Email = strings.TrimSpace(in.Email)
	in.PhoneNumber = strings.TrimSpace(in.PhoneNumber)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)

	if in.FirstName == "" || in.LastName == "" {
		return &domain.AppError{Err: domain.ErrInvalidInput, Message: "first name and last name are required"}
	}
	if !emailRegex.MatchString(in.Email) {
		return &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid email address"}
	}
	if !phoneRegex.MatchString(in.PhoneNumber) {
		return &domain.AppError{Err: domain.ErrInvalidInput, Message: "invalid phone number format (e.g. +1234567890)"}
	}
	return nil
}
