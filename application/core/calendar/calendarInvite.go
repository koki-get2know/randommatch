package calendar

import (
	"bytes"
	"fmt"
	"log"
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
	"github.com/koki/randommatch/database"
	"github.com/koki/randommatch/entity"
	"github.com/koki/randommatch/matcher"
	"gopkg.in/gomail.v2"
)

func generateIcsInvitation(sender, subject, description string, attendees []string, duration int64, startAt time.Time) ([]byte, error) {
	language := "en-us"
	location := "It's up to you"

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	event := cal.AddEvent(uuid.New().String())
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(startAt)
	event.SetEndAt(startAt.Add(time.Duration(duration) * time.Minute))
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
	log.Println(cal.Serialize())
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

	return emailRaw.Bytes(), err
}

func SendInvite(matches []matcher.Match, orgaUid string, duration int64, inviteDate time.Time) (string, error) {
	jobId := uuid.New().String()
	if err := database.CreateJobStatus(jobId); err != nil {
		return "", err
	}

	errorschannel := make(chan []string)
	statuschannel := make(chan string)

	go func() {
		defer close(statuschannel)
		defer close(errorschannel)
		statuschannel <- "Running"
		errors := []string{}
		for _, match := range matches {
			match := match
			if err := sendInvite(&match, duration, inviteDate); err != nil {
				uids := []string{}
				for _, user := range match.Users {
					uids = append(uids, user.Id)
				}
				errors = append(errors, fmt.Sprintf("users: %v desc: %v", strings.Join(uids, ", "), err.Error()))
			}
		}
		if len(errors) > 0 {
			errorschannel <- errors
			statuschannel <- "Failed"
		} else {
			errorschannel <- errors // To be able to have the errors in both cases (empty or not)
			statuschannel <- "Done"
		}
	}()
	go func() {
		for {
			select {
			case mailJobStatus := <-statuschannel:
				log.Println("received", mailJobStatus)
				database.UpdateJobStatus(jobId, database.JobStatus(mailJobStatus))
				if mailJobStatus == "Done" || mailJobStatus == "Failed" {
					return
				}
			case mailErrors := <-errorschannel:
				log.Println("received", mailErrors)
				if len(mailErrors) > 0 {
					database.UpdateJobErrors(jobId, mailErrors)
				}
				_, err := saveMatchingInfo(len(matches), len(mailErrors), orgaUid)

				if err != nil {
					mailErrors = append(mailErrors, err.Error())
					database.UpdateJobErrors(jobId, mailErrors)
				}
			}
		}
	}()
	return jobId, nil
}

func saveMatchingInfo(numGroups int, numFailed int, orgaUid string) (string, error) {
	numConversations := numGroups - numFailed

	matchingStat := entity.MatchingStat{
		NumGroups:        numGroups,
		NumConversations: numConversations,
		NumFailed:        numFailed,
	}

	return database.CreateMatchingStat(matchingStat, orgaUid)
}

func sendInvite(match *matcher.Match, duration int64, inviteDate time.Time) error {
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
	names := []string{}
	userIds := []string{}
	for _, user := range match.Users {
		names = append(names, user.Name)
		userIds = append(userIds, user.Id)
	}

	uidsMails, err := database.GetEmailsFromUIds(userIds)

	if err != nil {
		log.Println(err)
		return err
	}
	emails := []string{}
	for _, email := range uidsMails {
		emails = append(emails, email)
	}

	description := fmt.Sprintf(`Hello %v,
	You have been matched to connect during this cycle.
	
	Please accept (one of the) proposed invite(s).
	In case of conflict, contact your matching peer(s) to arrange the Koki conversation another time.
	
	Happy Koki!`, strings.Join(names, ", "))

	rawMessage, err := generateIcsInvitation(sender, subject, description, emails, duration, inviteDate)
	if err != nil {
		log.Println(err)
		return err
	}
	// Create a new session in the us-east-1 region.
	// Replace us-east-1 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("SES_KEY_ID"), os.Getenv("SES_KEY_SECRET"), ""),
	})
	if err != nil {
		log.Println(err)
		return err
	}
	//Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{Data: rawMessage},
		Source:     aws.String(sender),
	}

	_, err = svc.SendRawEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			log.Println(err.Error())
		}

		return err
	}
	return nil
}
