package queue

import (
	"crypto/rand"
	"encoding/json"
	"net/mail"
	"os"
	"time"

	"github.com/jhillyerd/enmime"
	"github.com/oklog/ulid"
)

func AddToQueue(queue QueueRef, envelope *enmime.Envelope) error {

	// ensure existence of queue directory
	err := os.MkdirAll(
		QueuePath(queue),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	// prepare message header
	var msg Message
	msg.Participants.To = []Participant{}
	msg.Attachments = []Attachment{}

	dateStr := envelope.GetHeader("Date")
	if len(dateStr) > 0 {
		formats := []string{
			"Mon, _2 Jan 2006 15:04:05 MST",
			"Mon, _2 Jan 2006 15:04:05 -0700",
			time.RFC1123,
			time.RFC1123Z,
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.RFC822,
			time.RFC822Z,
			time.RFC850,
			time.RFC3339,
			time.RFC3339Nano,
		}
		var timestamp time.Time
		var err error
		for _, format := range formats {
			timestamp, err = time.Parse(format, dateStr)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
		msg.Timestamp = timestamp
	} else {
		msg.Timestamp = time.Now()
	}

	addresses, err := envelope.AddressList("From")
	if err != mail.ErrHeaderNotPresent {
		if err != nil {
			return err
		}

		for _, address := range addresses {
			msg.Participants.From = &Participant{
				Name:  address.Name,
				Email: address.Address,
			}
			break
		}
	}

	addresses, err = envelope.AddressList("To")
	if err != mail.ErrHeaderNotPresent {
		if err != nil {
			return err
		}

		for _, address := range addresses {
			msg.Participants.To = append(
				msg.Participants.To,
				Participant{
					Name:  address.Name,
					Email: address.Address,
				},
			)
		}
	}

	msg.Subject = envelope.GetHeader("Subject")

	// prepare message body
	msg.Body.Text = envelope.Text
	msg.Body.HTML = envelope.HTML

	// add attachment data
	for _, attachment := range envelope.Attachments {
		msg.Attachments = append(
			msg.Attachments,
			Attachment{
				Id:      attachment.ContentID,
				Type:    attachment.ContentType,
				Charset: attachment.Charset,
				Name:    attachment.FileName,
				Content: attachment.Content,
			},
		)
	}

	// generate a ulid for message
	id, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		return err
	}

	// write message data
	pathPreparing := MessagePath(queue, id, MessageStatusPreparing)

	file, err := os.Create(pathPreparing)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)

	if err = encoder.Encode(msg); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	// change message state to queued
	pathQueued := MessagePath(queue, id, MessageStatusQueued)

	err = os.Rename(
		pathPreparing,
		pathQueued,
	)
	if err != nil {
		return err
	}

	return nil
}
