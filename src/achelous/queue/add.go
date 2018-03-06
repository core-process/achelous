package queue

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
)

func AddToQueue(envelope *enmime.Envelope) error {

	// generate a uuid for message
	newID, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	// create message directory
	err = os.MkdirAll(
		MsgBasePath(MessageStatusPreparing, newID),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	// prepare message meta data
	var msgMeta MessageMeta

	msgMeta.Timestamp = time.Now()
	msgMeta.Subject = envelope.GetHeader("Subject")

	addresses, err := envelope.AddressList("From")
	if err != nil {
		return err
	}

	for _, address := range addresses {
		msgMeta.Participants.From = &Participant{
			Name:  address.Name,
			Email: address.Address,
		}
		break
	}

	addresses, err = envelope.AddressList("To")
	if err != nil {
		return err
	}

	for _, address := range addresses {
		msgMeta.Participants.To = append(
			msgMeta.Participants.To,
			Participant{
				Name:  address.Name,
				Email: address.Address,
			},
		)
	}

	// write message meta data
	msgMetaJSON, err := json.Marshal(msgMeta)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(
		MsgMetaPath(MessageStatusPreparing, newID),
		msgMetaJSON,
		0644,
	)
	if err != nil {
		return err
	}

	// write text body
	if len(envelope.Text) > 0 {
		err = ioutil.WriteFile(
			MsgTextBodyPath(MessageStatusPreparing, newID),
			[]byte(envelope.Text),
			0644,
		)
		if err != nil {
			return err
		}
	}

	// write html body
	if len(envelope.HTML) > 0 {
		err = ioutil.WriteFile(
			MsgHtmlBodyPath(MessageStatusPreparing, newID),
			[]byte(envelope.HTML),
			0644,
		)
		if err != nil {
			return err
		}
	}

	// iterate attachments
	for _, attachment := range envelope.Attachments {

		// create directory
		err = os.MkdirAll(
			AttBasePath(MessageStatusPreparing, newID, attachment.ContentID),
			os.ModePerm,
		)
		if err != nil {
			return err
		}

		// prepare attachment meta
		var attMeta AttachmentMeta
		attMeta.Id = attachment.ContentID
		attMeta.Type = attachment.ContentType
		attMeta.Charset = attachment.Charset
		attMeta.Name = attachment.FileName

		// write attachment meta data
		attMetaJSON, err := json.Marshal(attMeta)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(
			AttMetaPath(MessageStatusPreparing, newID, attachment.ContentID),
			attMetaJSON,
			0644,
		)
		if err != nil {
			return err
		}

		// write attachment body
		err = ioutil.WriteFile(
			AttBodyPath(MessageStatusPreparing, newID, attachment.ContentID),
			attachment.Content,
			0644,
		)
		if err != nil {
			return err
		}
	}

	// create directory for queued messages
	err = os.MkdirAll(
		BasePath(MessageStatusQueued),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	// change message state to queued
	err = os.Rename(
		MsgBasePath(MessageStatusPreparing, newID),
		MsgBasePath(MessageStatusQueued, newID),
	)
	if err != nil {
		return err
	}

	return nil
}
