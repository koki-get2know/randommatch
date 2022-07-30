package calendar

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

func gomailing(subject, htmlBody string) ([]byte, error) {
	m := gomail.NewMessage()

	m.SetHeader("From", "Ivanov Tmib <brobizzness@gmail.com>")
	m.SetHeader("To", "tmibkage@yahoo.fr", "tmibkage@gmail.com", "ivan.tchomguemieguem@amadeus.com")
	m.SetAddressHeader("Cc", "brobizzness@gmail.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	m.Attach("Iterationreview.ics")
	var emailRaw bytes.Buffer
	_, err := m.WriteTo(&emailRaw)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("m: %v\n", emailRaw.String())
	return emailRaw.Bytes(), nil
}

func SendInvite() {

	const (
		// Replace sender@example.com with your "From" address.
		// This address must be verified with Amazon SES.
		Sender = "brobizzness@gmail.com"

		// Specify a configuration set. To use a configuration
		// set, comment the next line and line 92.
		//ConfigurationSet = "ConfigSet"

		// The subject line for the email.
		Subject = "Amazon SES Test (AWS SDK for Go)"

		// The HTML body for the email.
		HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
			"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
			"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"
	)

	rawMessage, err := gomailing(Subject, HtmlBody)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Create a new session in the us-east-1 region.
	// Replace us-east-1 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	if err != nil {
		fmt.Println(err)
		return
	}
	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{Data: rawMessage},
		Source:     aws.String(Sender),
	}

	// Attempt to send the email.
	//result, err := svc.SendEmail(input)
	result, err := svc.SendRawEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return
	}

	fmt.Println("Email Sent")
	fmt.Println(result)
	// Authentication.
	//auth := smtp.PlainAuth("ses-smtp-user.20220720-191343", "AKIAUQ5SMTNQKJH7CSR5", "BMiZx6u/2IC33GIt9fYQiz66c5+MC60fl/OixMatMCC/", smtpHost)
	//auth := smtp.PlainAuth("ses-smtp-user.20220720-191343", "AKIAUQ5SMTNQEKMPO3C5", "BMBwPqGwY88P/unnpJdtilHsFImWaQ7y+J2cuCT2vwGP", smtpHost)

}
