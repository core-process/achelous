package queue

import (
	"encoding/json"
	"os"
)

func Add(queue QueueRef, message Message, prettyJSON bool) error {

	// ensure existence of queue directory
	err := os.MkdirAll(
		QueuePath(queue),
		os.ModePerm,
	)
	if err != nil {
		return err
	}

	// write message data
	pathPreparing := MessagePath(queue, message.ID, MessageStatusPreparing)

	file, err := os.Create(pathPreparing)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if prettyJSON {
		encoder.SetIndent("", "  ")
	}
	encoder.SetEscapeHTML(false)

	if err = encoder.Encode(message); err != nil {
		return err
	}

	if err = file.Sync(); err != nil {
		return err
	}

	if err = file.Close(); err != nil {
		return err
	}

	// change message state to queued
	pathQueued := MessagePath(queue, message.ID, MessageStatusQueued)

	err = os.Rename(
		pathPreparing,
		pathQueued,
	)
	if err != nil {
		return err
	}

	return nil
}
