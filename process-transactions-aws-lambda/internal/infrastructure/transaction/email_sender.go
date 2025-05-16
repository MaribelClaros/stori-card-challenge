package transaction

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

const (
	copyRecipient = "claross.maribel@gmail.com"
	emailTemplate = `
	<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8" />
  <title>Stori Account Summary</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #ffffff;">
  <table width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; padding: 20px 0;">
    <tr>
      <td align="center">
        <table width="600" cellpadding="0" cellspacing="0" border="0" style="background-color: #ffffff; border: 1px solid #eee; border-radius: 10px; overflow: hidden;">
          <!-- Header -->
          <tr style="background-color: #00d180;">
            <td style="padding: 20px;" align="center">
              <img src="https://www.storicard.com/_next/static/media/stori_s_color.90dc745f.svg" alt="Stori Logo" width="120" />

            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding: 30px;">
              <h2 style="color: #00d180; margin-top: 0;">Your Monthly Account Summary</h2>
              <p style="font-size: 16px; color: #333;">Hello,</p>
              <p style="font-size: 16px; color: #333;">Here is your latest account summary with Stori:</p>

              <h3 style="color: #00d180; margin-bottom: 5px;">Total Balance:</h3>
              <p style="font-size: 26px; font-weight: bold; color: #333; margin: 10px 0;">$ {{.Balance}}</p>

              <h3 style="color: #00d180; margin-top: 30px;">Monthly Activity</h3>
              <table width="100%" cellpadding="8" cellspacing="0" border="0" style="border-collapse: collapse; margin-top: 10px;">
                <thead>
                  <tr style="background-color: #f5f5f5; color: #00d180;">
                    <th align="left">Month</th>
                    <th align="right">Transactions</th>
                    <th align="right">Average Debit</th>
                    <th align="right">Average Credit</th>
                  </tr>
                </thead>
                <tbody>
                  {{range .Monthly}}
                  <tr style="border-top: 1px solid #eee; color: #333;">
                    <td>{{.Month}}</td>
                    <td align="right">{{.TransactionCount}}</td>
                    <td align="right">${{printf "%.2f" .AverageDebit}}</td>
                    <td align="right">${{printf "%.2f" .AverageCredit}}</td>
                  </tr>
                  {{end}}
                </tbody>
              </table>

              <p style="font-size: 14px; color: #777; margin-top: 30px;">Thank you for using Stori.</p>
            </td>
          </tr>

          <!-- Footer -->
          <tr style="background-color: #f8f8f8;">
            <td align="center" style="padding: 20px; font-size: 12px; color: #999;">
              &copy; Stori. All rights reserved.
            </td>
          </tr>
        </table>
      </td>
    </tr>
  </table>
</body>
</html>

	`
	subject = "Stori Card Transactions Status"
	sender  = "claross.maribel@gmail.com"
)

type EmailSender interface {
	SendEmail(status *TransactionsStatus, recipient string) error
}

type emailSender struct {
	sesClient *ses.SES
}

func NewGetEmailSender(session *session.Session) *emailSender {
	return &emailSender{
		sesClient: ses.New(session),
	}
}

func (e *emailSender) SendEmail(status *TransactionsStatus, recipient string) error {

	emailContent, err := generateEmailContent(emailTemplate, status)
	if err != nil {
		log.Fatal("Error generating email content:", err)
	}
	toAddresses := []*string{
		aws.String(copyRecipient),
		aws.String(recipient),
	}
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: toAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(emailContent),
				},
			},
			Subject: &ses.Content{
				Data: aws.String(subject),
			},
		},
		Source: aws.String(sender),
	}

	// Send the email
	result, err := e.sesClient.SendEmail(input)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("Email sent successfully:", result)
	return nil
}

func generateEmailContent(templateStr string, data *TransactionsStatus) (string, error) {
	tmpl, err := template.New("email").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
