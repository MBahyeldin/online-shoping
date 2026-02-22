package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/online-cake-shop/backend/internal/domain"
)

type envelope map[string]any

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("write json", "error", err)
	}
}

func writeSuccess(w http.ResponseWriter, status int, data any) {
	writeJSON(w, status, envelope{"success": true, "data": data})
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	status, message := resolveError(err)
	writeJSON(w, status, envelope{
		"success": false,
		"error":   message,
	})
}

func resolveError(err error) (int, string) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		err = appErr.Unwrap()
	}

	msg := err.Error()
	if appErr != nil && appErr.Message != "" {
		msg = appErr.Message
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, msg
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, msg
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, msg
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, msg
	case errors.Is(err, domain.ErrInvalidInput):
		return http.StatusBadRequest, msg
	case errors.Is(err, domain.ErrOTPExpired),
		errors.Is(err, domain.ErrOTPInvalid),
		errors.Is(err, domain.ErrOTPAlreadyUsed):
		return http.StatusUnprocessableEntity, msg
	case errors.Is(err, domain.ErrRateLimitExceeded):
		return http.StatusTooManyRequests, msg
	case errors.Is(err, domain.ErrInsufficientStock):
		return http.StatusConflict, msg
	case errors.Is(err, domain.ErrEmptyCart):
		return http.StatusBadRequest, msg
	default:
		slog.Error("unhandled error", "error", err)
		return http.StatusInternalServerError, "an internal error occurred"
	}
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(io.LimitReader(r.Body, 1<<20)) // 1 MB limit
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}
