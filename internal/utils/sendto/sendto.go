package sendto

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"go_microservice_backend_api/global"
	"html/template"
	"net/smtp"
	"strings"
)

var (
	SMTPHost     = global.Config.SMTP.Host
	SMTPPort     = global.Config.SMTP.Port
	SMTPUsername = global.Config.SMTP.Username
	SMTPPassword = global.Config.SMTP.Password
)

type EmailAddress struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
type Mail struct {
	From    EmailAddress
	To      []string
	Subject string
	Body    string
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.From.Address)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func SendTextEmailOtp(to []string, from string, otp string) error {
	contentEmail := Mail{
		From:    EmailAddress{Address: from, Name: "test"},
		To:      to,
		Subject: "OTP verification",
		Body:    fmt.Sprintf("Your OTP is %s. Please enter it to verify your account.", otp),
	}

	messageMail := BuildMessage(contentEmail)

	// send email

	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)

	err := smtp.SendMail(SMTPHost+":587", auth, from, to, []byte(messageMail))
	if err != nil {
		global.Logger.Error("Email send fail::", zap.Error(err))
		return err
	}
	return nil
}

func SendTemplateEmailOTP(
	to []string,
	from string,
	nameTemplate string,
	dataTemplate map[string]interface{}) error {
	htmlBody, err := getMailTemplate(nameTemplate, dataTemplate)
	if err != nil {
		return err
	}

	return send(to, from, htmlBody)
}
func getMailTemplate(nameTemplate string, dataTemplate map[string]interface{}) (string, error) {
	htmlTemplate := new(bytes.Buffer)
	t := template.Must(template.New(nameTemplate).ParseFiles("templates-email/" + nameTemplate))
	err := t.Execute(htmlTemplate, dataTemplate)
	if err != nil {
		return "", err
	}
	return htmlTemplate.String(), nil
}
func send(to []string, from string, htmlTemplate string) error {
	contentEmail := Mail{
		From:    EmailAddress{Address: from, Name: "test"},
		To:      to,
		Subject: "OTP verification",
		Body:    htmlTemplate,
	}

	messageMail := BuildMessage(contentEmail)

	// send email

	auth := smtp.PlainAuth("", SMTPUsername, SMTPPassword, SMTPHost)

	err := smtp.SendMail(SMTPHost+":587", auth, from, to, []byte(messageMail))
	if err != nil {
		//global.Logger.Error("Email send fail::", zap.Error(err))
		return err
	}
	return nil
}
