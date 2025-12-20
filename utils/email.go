package utils

import (
	"fmt"
	"log"
)

// SendPasswordResetEmail simulates sending a password reset email.
// In a real application, this would use an email service provider (e.g., SendGrid, Mailgun).
func SendPasswordResetEmail(email, token string) error {
	// Construct the reset link
	resetLink := fmt.Sprintf("http://localhost:8080/reset-password-page?token=%s", token) // Example link

	// Simulate sending the email by printing to the console
	log.Println("========================================================")
	log.Printf("SIMULATING SENDING PASSWORD RESET EMAIL")
	log.Printf("To: %s", email)
	log.Printf("Subject: Reset Your Password")
	log.Printf("Body: To reset your password, please click the following link: %s", resetLink)
	log.Println("========================================================")

	return nil // In a real scenario, return any error from the email client
}
