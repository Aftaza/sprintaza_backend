package utils

import (
	"bytes"
	"html/template"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

// SendEmail sends an email using SendGrid.
func SendEmail(toName, toEmail, subject, templatePath string, data interface{}) error {
	// Get SendGrid API key from environment
	apiKey := GodotEnv("SENDGRID_API_KEY")
	if apiKey == "" {
		logrus.Error("SENDGRID_API_KEY not set in environment")
		return nil // Or return an error if email is critical
	}

	// Get sender name and email from environment
	fromName := GodotEnv("SENDGRID_FROM_NAME")
	fromEmail := GodotEnv("SENDGRID_FROM_EMAIL")

	// Create a new SendGrid client
	client := sendgrid.NewSendClient(apiKey)

	// Create a new email message
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	// Parse the HTML template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"template_path": templatePath,
			"error":         err.Error(),
		}).Error("Failed to parse email template")
		return err
	}

	// Execute the template with the provided data
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		logrus.WithFields(logrus.Fields{
			"template_path": templatePath,
			"error":         err.Error(),
		}).Error("Failed to execute email template")
		return err
	}

	// Create the email content
	message := mail.NewSingleEmail(from, subject, to, "", body.String())

	// Send the email
	response, err := client.Send(message)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"to_email": toEmail,
			"error":    err.Error(),
		}).Error("Failed to send email")
		return err
	}

	// Log the response
	if response.StatusCode >= 300 {
		logrus.WithFields(logrus.Fields{
			"to_email":    toEmail,
			"status_code": response.StatusCode,
			"body":        response.Body,
		}).Warn("SendGrid returned a non-success status code")
	} else {
		logrus.WithFields(logrus.Fields{
			"to_email": toEmail,
		}).Info("Email sent successfully")
	}

	return nil
}
