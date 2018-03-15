package programs

import (
	"bufio"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/core-process/achelous/common/queue"
	"github.com/core-process/achelous/spring-core/args"
	"github.com/core-process/achelous/spring-core/config"
	"github.com/core-process/achelous/spring-core/debug"

	"github.com/jhillyerd/enmime"
	"github.com/oklog/ulid"
)

func Sendmail(cdata *config.Config, smArgs *args.SmArgs, recipients []string) error {

	// prepare input
	ignoreDot := smArgs.Arg_i || smArgs.Arg_oi || smArgs.Arg_O.Opt_IgnoreDots

	filterReader, filterWriter := io.Pipe()

	debugInfoStoppedAtDotLine := false
	debugInfoEndOfHeaderSimulated := false

	go func() {
		// close writer at end of func
		defer filterWriter.Close()

		// scan lines, perform modifications
		scanner := bufio.NewScanner(os.Stdin)
		headerSection := true
		currentLine := -1

		for scanner.Scan() {
			// get next line
			line := scanner.Text()
			currentLine++
			// handle "."
			if !ignoreDot && line == "." {
				debugInfoStoppedAtDotLine = true
				break
			}
			// handle end of header
			if headerSection {
				if len(line) == 0 {
					// we leave the header section naturally
					headerSection = false
				} else {
					// we have to simulate the end of the header section
					isHeaderLine := strings.Index(line, ":") > 0
					isContLine := strings.IndexAny(line, " \t") == 0 && currentLine > 0

					if !isHeaderLine && !isContLine {
						debugInfoEndOfHeaderSimulated = true
						line = "\n" + line
						headerSection = false
					}
				}
			}
			// forward data (and errors)
			line += "\n"
			if _, err := io.WriteString(filterWriter, line); err != nil {
				filterWriter.CloseWithError(err)
				break
			}
		}
		// forward errors
		if err := scanner.Err(); err != nil {
			filterWriter.CloseWithError(err)
		}
	}()

	// parse envelope
	envelope, err := enmime.ReadEnvelope(filterReader)
	if err != nil {
		debug.MailOn(debug.OtherErrors, cdata, "parsing email", nil, err, struct {
			IgnoreDotEnabled     bool
			StoppedAtDotLine     bool
			EndOfHeaderSimulated bool
		}{
			ignoreDot,
			debugInfoStoppedAtDotLine,
			debugInfoEndOfHeaderSimulated,
		})
		return err
	}

	// check for parsing errors
	if len(envelope.Errors) > 0 {
		// create error
		errMsg := "Parsing failed:"
		for _, v := range envelope.Errors {
			errMsg += "\n- " + v.String()
		}
		err = errors.New(errMsg)
		// handle error
		debug.MailOn(debug.ParsingErrors, cdata, "parsing error validation", nil, err, envelope.Errors)
		if cdata.AbortOnParseErrors {
			return err
		}
		log.Printf("Ignoring parse errors: %v", err)
	}

	// feed values from args to from and to
	if len(envelope.Root.Header.Get("From")) == 0 {
		if smArgs.Arg_f != nil && smArgs.Arg_F != nil {
			envelope.Root.Header.Set("From", *smArgs.Arg_F+" <"+*smArgs.Arg_f+">")
		} else if smArgs.Arg_f != nil {
			envelope.Root.Header.Set("From", *smArgs.Arg_f)
		} else if smArgs.Arg_F != nil {
			envelope.Root.Header.Set("From", *smArgs.Arg_F)
		}
	}

	if len(envelope.Root.Header.Get("To")) == 0 && len(recipients) > 0 {
		envelope.Root.Header.Set("To", strings.Join(recipients, ", "))
	}

	// create message
	msgID, err := ulid.New(ulid.Now(), rand.Reader)
	if err != nil {
		debug.MailOn(debug.OtherErrors, cdata, "creating message id", nil, err, nil)
		return err
	}

	message := queue.Message{ID: msgID}

	message.Participants.To = []queue.Participant{}
	message.Attachments = []queue.Attachment{}

	// extract "timestamp"
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
			debug.MailOn(debug.ParsingErrors, cdata, "parsing timestamp", &msgID, err, dateStr)
			if cdata.AbortOnParseErrors {
				return err
			}
			log.Printf("Ignoring parse errors: %v", err)
			timestamp = time.Now()
		}
		message.Timestamp = timestamp
	} else {
		message.Timestamp = time.Now()
	}

	// extract "from"
	addresses, err := envelope.AddressList("From")
	if err != mail.ErrHeaderNotPresent {
		if err != nil && err.Error() == "mail: no angle-addr" {
			message.Participants.From = &queue.Participant{
				Name:  envelope.GetHeader("From"),
				Email: "",
			}
		} else if err != nil {
			debug.MailOn(debug.ParsingErrors, cdata, "parsing 'from' field", &msgID, err, envelope.GetHeader("From"))
			if cdata.AbortOnParseErrors {
				return err
			}
			log.Printf("Ignoring parse errors: %v", err)
			message.Participants.From = &queue.Participant{
				Name:  envelope.GetHeader("From"),
				Email: "",
			}
		} else {
			for _, address := range addresses {
				message.Participants.From = &queue.Participant{
					Name:  address.Name,
					Email: address.Address,
				}
				break
			}
		}
	}

	// extract "to"
	addresses, err = envelope.AddressList("To")
	if err != mail.ErrHeaderNotPresent {
		if err != nil && err.Error() == "mail: no angle-addr" {
			message.Participants.To = append(
				message.Participants.To,
				queue.Participant{
					Name:  envelope.GetHeader("To"),
					Email: "",
				},
			)
		} else if err != nil {
			debug.MailOn(debug.ParsingErrors, cdata, "parsing 'to' field", &msgID, err, envelope.GetHeader("To"))
			if cdata.AbortOnParseErrors {
				return err
			}
			log.Printf("Ignoring parse errors: %v", err)
			message.Participants.To = append(
				message.Participants.To,
				queue.Participant{
					Name:  envelope.GetHeader("To"),
					Email: "",
				},
			)
		} else {
			for _, address := range addresses {
				message.Participants.To = append(
					message.Participants.To,
					queue.Participant{
						Name:  address.Name,
						Email: address.Address,
					},
				)
			}
		}
	}

	// extract "subject"
	message.Subject = envelope.GetHeader("Subject")

	// extract "body"
	message.Body.Text = envelope.Text
	message.Body.HTML = envelope.HTML

	// extract "attachments"
	for _, attachment := range envelope.Attachments {
		message.Attachments = append(
			message.Attachments,
			queue.Attachment{
				ID:      attachment.ContentID,
				Type:    attachment.ContentType,
				Charset: attachment.Charset,
				Name:    attachment.FileName,
				Content: attachment.Content,
			},
		)
	}

	// add message to queue
	err = queue.Add(
		queue.QueueRef(cdata.DefaultQueue),
		message,
		cdata.PrettyJSON,
	)
	if err != nil {
		debug.MailOn(debug.OtherErrors, cdata, "submitting email to queue", &msgID, err, message)
		return err
	}

	// trigger queue run
	if cdata.TriggerQueueRun {
		err = queue.Trigger()
		if err != nil {
			log.Printf("failed to trigger queue run: %v", err)
		}
	}

	return nil
}
