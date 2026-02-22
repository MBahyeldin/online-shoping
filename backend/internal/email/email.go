package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/online-cake-shop/backend/internal/config"
)

// Sender is the abstract email provider interface.
type Sender interface {
	SendOTP(to, firstName, otp string) error
}

// ─── SMTP Implementation ──────────────────────────────────────────────────────

type SMTPSender struct {
	cfg config.EmailConfig
}

func NewSMTPSender(cfg config.EmailConfig) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) SendOTP(to, firstName, otp string) error {
	subject := "Your Cake Shop Verification Code"
	body, err := renderOTPTemplate(firstName, otp)
	if err != nil {
		return fmt.Errorf("render otp template: %w", err)
	}

	msg := buildMIMEMessage(s.cfg.From, to, subject, body)

	auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPass, s.cfg.SMTPHost)
	addr := fmt.Sprintf("%s:%d", s.cfg.SMTPHost, s.cfg.SMTPPort)

	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg))
}

func buildMIMEMessage(from, to, subject, htmlBody string) string {
	return fmt.Sprintf(
		"From: Cake Shop <%s>\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, to, subject, htmlBody,
	)
}

// ─── OTP Email Template ───────────────────────────────────────────────────────

const otpEmailTpl = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
  <title>Your Verification Code</title>
  <style>
    body { font-family: Arial, sans-serif; background: #f9f9f9; margin: 0; padding: 0; }
    .container { max-width: 480px; margin: 40px auto; background: #fff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,.08); }
    .header { background: #c05621; padding: 28px 32px; text-align: center; }
    .header h1 { color: #fff; margin: 0; font-size: 22px; letter-spacing: .5px; }
    .body { padding: 32px; }
    .body p { color: #444; line-height: 1.6; }
    .otp { font-size: 36px; font-weight: 700; letter-spacing: 8px; color: #c05621; text-align: center; margin: 24px 0; }
    .footer { background: #f3f3f3; padding: 16px 32px; text-align: center; font-size: 12px; color: #888; }
  </style>
</head>
<body>
  <div class="container">
    <div class="header"><h1>🎂 Cake Shop</h1></div>
    <div class="body">
      <p>Hello <strong>{{.FirstName}}</strong>,</p>
      <p>Use the verification code below to complete your registration. This code expires in <strong>5 minutes</strong>.</p>
      <div class="otp">{{.OTP}}</div>
      <p>If you didn't request this, you can safely ignore this email.</p>
    </div>
    <div class="footer">© 2024 Cake Shop. All rights reserved.</div>
  </div>
</body>
</html>`

type otpTemplateData struct {
	FirstName string
	OTP       string
}

func renderOTPTemplate(firstName, otp string) (string, error) {
	tpl, err := template.New("otp").Parse(otpEmailTpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, otpTemplateData{FirstName: firstName, OTP: otp}); err != nil {
		return "", err
	}
	return buf.String(), nil
}
