package calendar

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/google/uuid"
	"github.com/koki/randommatch/matcher"
	"gopkg.in/gomail.v2"
)

func generateIcsInvitation(sender, subject, description string, attendees []string) ([]byte, error) {
	startAt := time.Now()
	language := "en-us"
	location := "It's up to you"

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	event := cal.AddEvent(uuid.New().String())
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(startAt)
	event.SetEndAt(startAt.Add(30 * time.Minute))
	event.SetSummary(subject, &ics.KeyValues{Key: string(ics.ParameterLanguage), Value: []string{language}})
	event.SetLocation(location)
	event.SetDescription(description)
	event.SetOrganizer("mailto:"+sender, ics.WithCN("Koki Admin"))
	event.SetTimeTransparency(ics.TransparencyOpaque)

	for _, attendee := range attendees {
		event.AddAttendee(attendee, ics.CalendarUserTypeIndividual, ics.ParticipationStatusNeedsAction, ics.ParticipationRoleReqParticipant, ics.WithRSVP(true))
	}
	event.SetSequence(1)
	reminder := event.AddAlarm()
	reminder.SetAction(ics.ActionDisplay)
	reminder.SetTrigger("-P15M")
	fmt.Println(cal.Serialize())
	filename := "invite.ics"
	m := gomail.NewMessage()

	m.SetHeader("subject", subject)
	m.SetHeader("From", fmt.Sprintf("Koki Admin <%v>", sender))
	m.SetHeader("To", strings.Join(attendees, ", "))
	m.SetHeader("Content-Description", filename)
	m.SetHeader("Content-class", "urn:content-classes:calendarmessage")
	m.SetHeader("Filename", filename)
	m.SetHeader("Path", filename)

	m.SetBody("text/plain", description)
	m.AddAlternative(`text/calendar; method="REQUEST"; name="invite.ics"`,
		cal.Serialize(),
		gomail.SetPartEncoding(gomail.Base64))

	var emailRaw bytes.Buffer
	_, err := m.WriteTo(&emailRaw)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("m: %v\n", emailRaw.String())
	return emailRaw.Bytes(), nil
}

func SendInvite(match *matcher.Match, adminEmail string) error {

	const (
		// Replace sender@example.com with your "From" address.
		// This address must be verified with Amazon SES.
		sender = "koki.get2know@gmail.com"

		// Specify a configuration set. To use a configuration
		// set, comment the next line and line 92.
		//ConfigurationSet = "ConfigSet"

		// The subject line for the email.
		subject = "You are matched for a Koki!"
	)
	userIds := []string{}
	for _, user := range match.Users {
		userIds = append(userIds, user.Id)
	}

	description := fmt.Sprintf(`Hello %v,
	You have been matched to connect during this cycle.
	
	Please accept (one of the) proposed invite(s).
	In case of conflict, contact your matching peer(s) to arrange the Koki conversation another time.
	
	Happy Koki!`, strings.Join(userIds, ", "))

	attendees := []string{adminEmail}

	rawMessage, err := generateIcsInvitation(sender, subject, description, attendees)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Create a new session in the us-east-1 region.
	// Replace us-east-1 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("SES_KEY_ID"), os.Getenv("SES_KEY_SECRET"), ""),
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{Data: rawMessage},
		Source:     aws.String(sender),
	}

	result, err := svc.SendRawEmail(input)

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

		return err
	}

	fmt.Println("Email Sent", result)
	return nil
}
