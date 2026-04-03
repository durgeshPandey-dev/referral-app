package email

import (
	"context"
	"fmt"
	"os"

	"github.com/durgeshPandey-dev/referral/apps/api/logger"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func Send(ctx context.Context, to, subject, htmlContent string) error {
	from := mail.NewEmail("Referral App", os.Getenv("SENDER_EMAIL"))
	toEmail := mail.NewEmail("", to)

	// Plain text fallback (important for production)
	plainText := "Please view this email in HTML format."

	message := mail.NewSingleEmail(
		from,
		subject,
		toEmail,
		plainText,
		htmlContent,
	)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	resp, err := client.Send(message)
	if err != nil {
		logger.Error(ctx, "sendgrid_request_failed", map[string]interface{}{
			"email": to,
			"error": err,
		})
		return err
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("sendgrid error: status=%d body=%s", resp.StatusCode, resp.Body)
		logger.Error(ctx, "sendgrid_rejected", map[string]interface{}{
			"email":       to,
			"status_code": resp.StatusCode,
		})
		return err
	}

	logger.Info(ctx, "email_sent_success", map[string]interface{}{
		"email":       to,
		"status_code": resp.StatusCode,
	})

	return nil
}
