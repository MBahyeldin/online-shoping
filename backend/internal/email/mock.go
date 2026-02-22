package email

import "log/slog"

// MockSender logs emails to stdout — for development and testing.
type MockSender struct {
	logger *slog.Logger
}

func NewMockSender(logger *slog.Logger) *MockSender {
	return &MockSender{logger: logger}
}

func (m *MockSender) SendOTP(to, firstName, otp string) error {
	m.logger.Info("📧 [MOCK EMAIL] OTP sent",
		"to", to,
		"firstName", firstName,
		"otp", otp,
	)
	return nil
}
