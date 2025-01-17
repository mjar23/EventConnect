// email.go
package emailUtil

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func SendEmail(to []string, subject, plainTextContent, htmlContent string) error {
	// Get the SendGrid API key from an environment variable
	sendGridAPIKey := os.Getenv("SENDGRID_API_KEY")
	if sendGridAPIKey == "" {
		return fmt.Errorf("SendGrid API key not set in environment variables")
	}

	from := mail.NewEmail("Event-Connect Team", "mauricejarvis80@gmail.com")
	toEmail := mail.NewEmail("", to[0]) // Create an Email object for the first recipient

	message := mail.NewSingleEmail(from, subject, toEmail, plainTextContent, htmlContent)

	// Add additional recipients if there are any
	for _, recipient := range to[1:] {
		message.Personalizations[0].AddTos(mail.NewEmail("", recipient))
	}

	client := sendgrid.NewSendClient(sendGridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	fmt.Printf("Email sent with status code: %d\n", response.StatusCode)
	fmt.Printf("Response body: %s\n", response.Body)
	fmt.Printf("Response headers: %v\n", response.Headers)

	return nil
}
