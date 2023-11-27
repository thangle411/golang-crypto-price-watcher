package email

import (
	"net/smtp"

	"github.com/jordan-wright/email"
)

type Email struct {
	AppEmail    string
	AppPassword string
	ToEmail     []string
}

func SendEmail(subject string, htmlString string, emailConfig Email) error {
	e := email.NewEmail()
	e.From = "Price Tracker <" + emailConfig.AppEmail + ">"
	e.To = emailConfig.ToEmail
	e.Subject = subject + " price changed"
	e.HTML = []byte(htmlString)
	return e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "golanglearner411@gmail.com", emailConfig.AppPassword, "smtp.gmail.com"))
}

func TestEmail(receiverEmail string, senderEmail string, appPassword string) error {
	err := SendEmail("test ticker", `<div>test</div>`, Email{
		AppEmail:    senderEmail,
		AppPassword: appPassword,
		ToEmail:     []string{receiverEmail},
	})
	if err != nil {
		return err
	}
	return nil
}
