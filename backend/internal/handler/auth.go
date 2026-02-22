package handler

import (
	"net/http"
	"time"

	"github.com/online-cake-shop/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// ─── Register ────────────────────────────────────────────────────────────────

type registerRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"success": false, "error": "invalid request body"})
		return
	}

	if err := h.authSvc.Register(r.Context(), service.RegisterInput{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}); err != nil {
		writeError(w, r, err)
		return
	}

	writeSuccess(w, http.StatusCreated, envelope{
		"message": "OTP sent to your email address. Please verify to complete registration.",
	})
}

// ─── Verify OTP ──────────────────────────────────────────────────────────────

type verifyOTPRequest struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}

func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req verifyOTPRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"success": false, "error": "invalid request body"})
		return
	}

	result, err := h.authSvc.VerifyOTP(r.Context(), service.VerifyOTPInput{
		Email: req.Email,
		OTP:   req.OTP,
	})
	if err != nil {
		writeError(w, r, err)
		return
	}

	// Set HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set true in production with TLS
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	writeSuccess(w, http.StatusOK, envelope{
		"token": result.Token,
		"user": envelope{
			"id":         result.User.ID.String(),
			"first_name": result.User.FirstName,
			"last_name":  result.User.LastName,
			"email":      result.User.EmailAddress,
			"phone":      result.User.PhoneNumber,
		},
	})
}

// ─── Resend OTP ──────────────────────────────────────────────────────────────

type resendOTPRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) ResendOTP(w http.ResponseWriter, r *http.Request) {
	var req resendOTPRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, envelope{"success": false, "error": "invalid request body"})
		return
	}

	if err := h.authSvc.ResendOTP(r.Context(), req.Email); err != nil {
		writeError(w, r, err)
		return
	}

	writeSuccess(w, http.StatusOK, envelope{
		"message": "A new OTP has been sent to your email address.",
	})
}
