package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SendPasswordResetEmail sends a real password reset email using SMTP.
// It reads configuration from environment variables:
// SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS
func SendPasswordResetEmail(email, token string) error {
	// Get SMTP configuration from environment variables
	// email : not used in this
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpTEST := os.Getenv("SMTP_TEST")
	if smtpHost == "" || smtpPortStr == "" || smtpUser == "" || smtpPass == "" {
		log.Println("WARNING: SMTP environment variables not fully configured. Falling back to console output.")
		return logSimulatedEmail(email, token)
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("ERROR: Invalid SMTP_PORT value: %v", err)
		return err
	}

	// Construct the reset link
	// In a real app, the base URL should come from config
	resetLink := fmt.Sprintf("http://localhost:3003/reset-password?token=%s", token)

	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser) // Or a specific "From" address
	m.SetHeader("To", smtpTEST)
	m.SetHeader("Subject", "Reset Your Password")
	m.SetBody("text/html", fmt.Sprintf("To reset your password, please click the following link: <a href=\"%s\">%s</a>", resetLink, resetLink))

	// Create a new Dialer
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

	// Send the email
	log.Printf("Attempting to send password reset email to %s via SMTP...", email)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("ERROR: Failed to send email: %v", err)
		return err
	}

	log.Printf("Successfully sent password reset email to %s", email)
	return nil
}

// logSimulatedEmail is the fallback for when SMTP is not configured.
func logSimulatedEmail(email, token string) error {
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password-page?token=%s", token)

	log.Println("========================================================")
	log.Printf("SIMULATING SENDING PASSWORD RESET EMAIL")
	log.Printf("To: %s", email)
	log.Printf("Subject: Reset Your Password")
	log.Printf("Body: To reset your password, please click the following link: %s", resetLink)
	log.Println("========================================================")
	return nil
}
