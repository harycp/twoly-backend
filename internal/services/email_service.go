package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type EmailService interface {
	SendPasswordResetOTP(toEmail, toName, otpCode string) error
}

type brevoEmailService struct{}

func NewEmailService() EmailService {
	return &brevoEmailService{}
}

func (s *brevoEmailService) SendPasswordResetOTP(toEmail, toName, otpCode string) error {
	apiKey := os.Getenv("BREVO_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("BREVO_API_KEY is not set")
	}

	senderEmail := os.Getenv("BREVO_SENDER_EMAIL")
	if senderEmail == "" {
		senderEmail = "noreply@twoly.com"
	}

	senderName := os.Getenv("BREVO_SENDER_NAME")
	if senderName == "" {
		senderName = "Twoly"
	}

	url := "https://api.brevo.com/v3/smtp/email"

	htmlContent := fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
	</head>
	<body style="margin: 0; padding: 0; font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background-color: #FFF7ED; color: #111827;">
		<table width="100%%" cellpadding="0" cellspacing="0" style="background-color: #FFF7ED; padding: 40px 20px;">
			<tr>
				<td align="center">
					<table width="100%%" max-width="500px" cellpadding="0" cellspacing="0" style="background-color: #ffffff; border-radius: 32px; padding: 40px; box-shadow: 0 10px 25px -5px rgba(248,180,200, 0.2); border: 1px solid #fdf2f8; text-align: center;">
						<tr>
							<td>
								<h1 style="color: #FDA4AF; font-size: 26px; font-weight: 900; letter-spacing: 3px; text-transform: uppercase; margin: 0 0 20px 0;">Twoly</h1>
								<h2 style="font-size: 20px; font-weight: 800; margin: 0 0 10px 0;">Password Reset Request</h2>
								<p style="color: #6b7280; font-size: 14px; line-height: 1.6; margin: 0 0 30px 0;">
									Hi %s,<br>We received a request to reset your password for your Twoly shared space. Use the OTP code below to securely authenticate.
								</p>
								<div style="background-color: #FFF7ED; border-radius: 20px; padding: 20px; margin: 0 0 30px 0; letter-spacing: 12px; font-size: 38px; font-weight: 900; color: #FDA4AF; border: 2px dashed #F8B4C8;">
									%s
								</div>
								<p style="color: #9ca3af; font-size: 12px; margin: 0; line-height: 1.5;">
									This code will expire in 15 minutes.<br>If you didn't request this, you can safely ignore this email.
								</p>
							</td>
						</tr>
					</table>
					<p style="color: #d1d5db; font-size: 12px; margin-top: 20px;">© 2026 Twoly. All rights reserved.</p>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, toName, otpCode)

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  senderName,
			"email": senderEmail,
		},
		"to": []map[string]string{
			{"email": toEmail, "name": toName},
		},
		"subject":     "Twoly - Your Password Reset Code",
		"htmlContent": htmlContent,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("api-key", apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, brevo status code: %d", res.StatusCode)
	}

	return nil
}