import email
import email.mime.multipart
import email.mime.text
import smtplib
import datetime as dt
import uuid
import icalendar
import pytz

# Imagine this function is part of a class which provides the necessary config data
def send_appointment(date, attendee_email, organiser_email, subj, description, location, start_hour, start_minute):
  # Timezone to use for our dates - change as needed
  tz = pytz.timezone("Europe/London")
  start = tz.localize(dt.datetime.combine(date, dt.time(start_hour, start_minute, 0)))
  # Build the event itself
  cal = icalendar.Calendar()
  cal.add('prodid', '-//My calendar application//example.com//')
  cal.add('version', '2.0')
  cal.add('method', "REQUEST")
  event = icalendar.Event()
  event.add('attendee', attendee_email)
  event.add('organizer', organiser_email)
  #event.add('status', "confirmed")
  #event.add('category', "Event")
  event.add('summary', subj)
  event.add('description', description)
  event.add('location', location)
  event.add('dtstart', start)
  event.add('dtend', tz.localize(dt.datetime.combine(date, dt.time(start_hour + 1, start_minute, 0))))
  event.add('dtstamp', tz.localize(dt.datetime.combine(date, dt.time(6, 0, 0))))
  event['uid'] = str(uuid.uuid4()) # Generate some unique ID
  #event.add('priority', 5)
  event.add('sequence', 1)
  event.add('created', tz.localize(dt.datetime.now()))

  # Add a reminder
  alarm = icalendar.Alarm()
  alarm.add("action", "DISPLAY")
  alarm.add('description', "Reminder")
  # The only way to convince Outlook to do it correctly
  #alarm.add("TRIGGER;RELATED=START", "-PT{0}H".format(9))
  alarm.add("TRIGGER;", "-P{0}M".format(15))

  #alarm.add("TRIGGER", "-P{0}M".format(15))

  event.add_component(alarm)
  cal.add_component(event)

  # Build the email message and attach the event to it
  msg = email.mime.multipart.MIMEMultipart("alternative")

  msg["Subject"] = subj
  msg["From"] = organiser_email
  msg["To"] = attendee_email
  msg["Content-class"] = "urn:content-classes:calendarmessage"

  msg.attach(email.mime.text.MIMEText(description))

  filename = "invite.ics"
  part = email.mime.base.MIMEBase('text', "calendar", method="REQUEST", name=filename)
  part.set_payload( cal.to_ical() )
  email.encoders.encode_base64(part)
  part.add_header('Content-Description', filename)
  part.add_header("Content-class", "urn:content-classes:calendarmessage")
  part.add_header("Filename", filename)
  part.add_header("Path", filename)
  msg.attach(part)

  print(msg)

  # Send the email out
  server = smtplib.SMTP('email-smtp.us-east-1.amazonaws.com', 587)
  server.starttls()
  server.login("AKIAUQ5SMTNQKJH7CSR5", "BMiZx6u/2IC33GIt9fYQiz66c5+MC60fl/OixMatMCC/")
  server.sendmail(msg["From"], [msg["To"]], msg.as_string())
  server.quit()

date = dt.datetime(2022, 8, 2)
location = "It's up to you"
sender = "koki.get2know@gmail.com"
subject = "You are matched for a Koki!"
description = """Hello name1, name2,
	You have been matched to connect during this cycle.
	
	Please accept (one of the) proposed invite(s).
	In case of conflict, contact your matching peer(s) to arrange the Koki conversation another time.
	
	Happy Koki!"""
send_appointment(date, "ivan.tchomguemieguem@amadeus.com", sender, subject, description, 
                location, 12, 25)