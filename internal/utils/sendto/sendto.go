package sendto

import (
	"bytes"
	"fmt"
	"go_microservice_backend_api/global"
	_const "go_microservice_backend_api/internal/const"
	gomail "gopkg.in/mail.v2"
	"html/template"
	"net/smtp"
	"strings"
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
	//contentEmail := Mail{
	//	From:    EmailAddress{Address: from, Name: "test"},
	//	To:      to,
	//	Subject: "OTP verification",
	//	Body:    fmt.Sprintf("Your OTP is %s. Please enter it to verify your account.", otp),
	//}
	//
	//messageMail := BuildMessage(contentEmail)
	//
	//// send email
	//auth := smtp.PlainAuth("", global.Config.SMTP.Username, global.Config.SMTP.Password, global.Config.SMTP.Host)
	//
	//err := smtp.SendMail(global.Config.SMTP.Host+":587", auth, from, to, []byte(messageMail))
	//if err != nil {
	//	global.Logger.Error("Email send fail::", zap.Error(err))
	//	return err
	//}
	// Create a new message
	message := gomail.NewMessage()

	// Set email headers
	message.SetHeader("From", _const.HOST_EMAIL)
	message.SetHeader("To", to[0])
	message.SetHeader("Subject", "Verify OTP")

	// Set email body
	message.SetBody("text/plain", fmt.Sprintf("Your OTP is %s. Please enter it to verify your account.", otp))

	// Set up the SMTP dialer
	dialer := gomail.NewDialer(global.Config.MailTrap.Host, global.Config.MailTrap.Port, global.Config.MailTrap.Username, global.Config.MailTrap.Password)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		fmt.Println("Error Send mail fail:: ", err)
		panic(err)
	} else {
		fmt.Println("Email sent successfully!")
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

	auth := smtp.PlainAuth("", global.Config.SMTP.Username, global.Config.SMTP.Password, global.Config.SMTP.Host)

	err := smtp.SendMail(global.Config.SMTP.Host+":587", auth, from, to, []byte(messageMail))
	if err != nil {
		//global.Logger.Error("Email send fail::", zap.Error(err))
		return err
	}
	return nil
}
