package sendto

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type MainRequest struct {
	ToEmail     string `json:"toEmail"`
	MessageBody string `json:"messageBody"`
	Subject     string `json:"subject"`
	Attachment  string `json:"attachment"`
}

func SendEmailToJavaByAPI(otp string, email string, purpose string) error {
	//URL
	postURL := "https://localhost:8080/email/send_text"

	mailRequest := MainRequest{
		ToEmail:     email,
		MessageBody: "OTP " + otp,
		Subject:     "Verify OTP" + purpose,
		Attachment:  "path/to/email",
	}

	requestBody, err := json.Marshal(mailRequest)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", postURL, bytes.NewBuffer(requestBody))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	return nil
}
