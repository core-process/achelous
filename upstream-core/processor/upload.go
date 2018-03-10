package processor

import (
	"log"

	commonQueue "github.com/core-process/achelous/common/queue"

	"github.com/oklog/ulid"
)

func upload(queue commonQueue.QueueRef, id ulid.ULID) error {
	log.Printf("fake-uploading message %s of queue /%s", id, queue)
	return nil
}
