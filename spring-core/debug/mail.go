package debug

import (
	"crypto/rand"
	"fmt"
	"log"
	"time"

	"github.com/core-process/achelous/common/queue"
	"github.com/core-process/achelous/spring-core/config"

	"github.com/oklog/ulid"
)

type ErrorCategory int8

const (
	InvalidParameters ErrorCategory = 1
	ParsingErrors     ErrorCategory = 2
	OtherErrors       ErrorCategory = 3
)

func MailOn(category ErrorCategory, cdata *config.Config, activity string, ref *ulid.ULID, activityErr error, data interface{}) {

	// check if we should generate a debug mail
	shouldSend := (category == InvalidParameters && cdata.GenerateDebugMail.OnInvalidParameters) ||
		(category == ParsingErrors && cdata.GenerateDebugMail.OnParsingErrors) ||
		(category == OtherErrors && cdata.GenerateDebugMail.OnOtherErrors)

	if !shouldSend {
		return
	}

	// create message
	msgID, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		log.Printf("failed to generate ulid for debug message: %v", err)
		return
	}

	message := queue.Message{ID: msgID}

	message.Participants.To = []queue.Participant{}
	message.Attachments = []queue.Attachment{}

	// set meta
	message.Timestamp = time.Now()
	message.Participants.From = &cdata.GenerateDebugMail.Message.Sender
	message.Participants.To = append(message.Participants.To, cdata.GenerateDebugMail.Message.Receiver)

	// set subject and body
	refStr := "None"
	if ref != nil {
		refStr = ref.String()
	}

	message.Subject = cdata.GenerateDebugMail.Message.Subject
	message.Body.Text = fmt.Sprintf(cdata.GenerateDebugMail.Message.Body, activity, refStr, activityErr, data)

	// add message to queue
	err = queue.Add(
		queue.QueueRef(cdata.DefaultQueue),
		message,
		cdata.PrettyJSON,
	)
	if err != nil {
		log.Printf("failed to add debug message to queue: %v", err)
		return
	}

	// trigger queue run
	if cdata.TriggerQueueRun {
		err = queue.Trigger()
		if err != nil {
			log.Printf("failed to trigger queue run: %v", err)
		}
	}
}
