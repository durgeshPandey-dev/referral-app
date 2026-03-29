package email

import (
	"context"
	"fmt"
	"os"
	"time"

	"referral-app/internal/observability"
	"referral-app/pkg/logger"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func Send(ctx context.Context, to, subject, htmlContent string) error {
	tracer := otel.Tracer("beckn-onix/email")
	ctx, span := tracer.Start(ctx, "email.send")
	defer span.End()
	span.SetAttributes(
		attribute.String("email.to", to),
		attribute.String("email.subject", subject),
	)

	start := time.Now()
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
		observability.ObserveEmailDuration(time.Since(start), "failed")
		logger.Error(ctx, "sendgrid_request_failed", map[string]interface{}{
			"email": to,
			"error": err,
		})
		span.RecordError(err)
		span.SetStatus(codes.Error, "sendgrid request failed")
		return err
	}

	if resp.StatusCode >= 400 {
		err = fmt.Errorf("sendgrid error: status=%d body=%s", resp.StatusCode, resp.Body)
		observability.ObserveEmailDuration(time.Since(start), "failed")
		logger.Error(ctx, "sendgrid_rejected", map[string]interface{}{
			"email":       to,
			"status_code": resp.StatusCode,
		})
		span.RecordError(err)
		span.SetStatus(codes.Error, "sendgrid rejected request")
		return err
	}

	observability.ObserveEmailDuration(time.Since(start), "success")
	logger.Info(ctx, "email_sent_success", map[string]interface{}{
		"email":       to,
		"status_code": resp.StatusCode,
	})
	span.SetAttributes(attribute.Int("sendgrid.status_code", resp.StatusCode))
	span.SetStatus(codes.Ok, "email sent")

	return nil
}
