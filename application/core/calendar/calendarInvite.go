package calendar

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/koki/randommatch/matcher"
	"gopkg.in/gomail.v2"
)

func gomailing(sender, subject, htmlBody string) ([]byte, error) {
	m := gomail.NewMessage()

	m.SetHeader("From", fmt.Sprintf("Koki Admin <%v>", sender))
	m.SetHeader("To", "tmibkage@gmail.com")
	m.SetAddressHeader("Cc", "brobizzness@gmail.com", "Dan")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)
	m.Attach("calendar/Iterationreview.ics")
	var emailRaw bytes.Buffer
	_, err := m.WriteTo(&emailRaw)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("m: %v\n", emailRaw.String())
	return emailRaw.Bytes(), nil
}

func SendInvite(match *matcher.Match) {

	const (
		// Replace sender@example.com with your "From" address.
		// This address must be verified with Amazon SES.
		sender = "brobizzness@gmail.com"

		// Specify a configuration set. To use a configuration
		// set, comment the next line and line 92.
		//ConfigurationSet = "ConfigSet"

		// The subject line for the email.
		subject = "You are matched for a Koki!"
	)
	userIds := []string{}
	for _, user := range match.Users {
		userIds = append(userIds, user.UserId)
	}
	messageHtmlBody := fmt.Sprintf(`<p>Hello %v<br>
	You have been matched to connect during this cycle.
	</p>
	<p>
	Please accept (one of the) proposed invite(s). In case of conflict, contact your matching peer(s) to arrange the Koki conversation another time.
	</p>
	Happy Koki!`, strings.Join(userIds, ", "))

	rawMessage, err := gomailing(sender, subject, messageHtmlBody)
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
		Source:     aws.String(sender),
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
