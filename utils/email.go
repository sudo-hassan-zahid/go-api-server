package utils

import (
	"fmt"
	"net/smtp"

	"github.com/sudo-hassan-zahid/go-api-server/internal/config"
)

func SendEmail(to []string, subject, body string, cfg config.SMTPConfig) error {
	var auth smtp.Auth
	if cfg.User != "" && cfg.Password != "" {
		auth = smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
	}

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+
		"%s\r\n", to[0], subject, body))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	if err := smtp.SendMail(addr, auth, cfg.From, to, msg); err != nil {
		return err
	}
	return nil
}

const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
<style>
body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
.container { max-width: 600px; margin: 20px auto; background: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); }
.header { text-align: center; padding-bottom: 20px; border-bottom: 1px solid #eee; }
.content { padding: 20px 0; color: #333; line-height: 1.6; }
.button { display: inline-block; padding: 10px 20px; background-color: #007bff; color: #ffffff; text-decoration: none; border-radius: 5px; margin-top: 10px; }
.footer { text-align: center; color: #999; font-size: 12px; margin-top: 20px; }
</style>
</head>
<body>
<div class="container">
	<div class="header">
		<h2>%s</h2>
	</div>
	<div class="content">
		%s
	</div>
	<div class="footer">
		&copy; 2024 Go API Server. All rights reserved.
	</div>
</div>
</body>
</html>
`

func SendVerificationEmail(email, token string, cfg config.SMTPConfig) error {
	link := fmt.Sprintf("http://localhost:8080/api/auth/verify-email?token=%s", token)
	content := fmt.Sprintf("<p>Thank you for registering. Please click the button below to verify your email address:</p><a href=\"%s\" class=\"button\">Verify Email</a>", link)
	body := fmt.Sprintf(emailTemplate, "Verify your email", content)
	return SendEmail([]string{email}, "Verify your email", body, cfg)
}

func SendPasswordResetEmail(email, token string, cfg config.SMTPConfig) error {
	content := fmt.Sprintf("<p>You requested a password reset. Here is your verification code:</p><h1 style=\"text-align: center; letter-spacing: 5px;\">%s</h1><p>This code will expire in 15 minutes.</p>", token)
	body := fmt.Sprintf(emailTemplate, "Reset Password Request", content)
	return SendEmail([]string{email}, "Reset your password", body, cfg)
}
