package processor

import (
	"log"

	commonQueue "github.com/core-process/achelous/common/queue"
	"github.com/core-process/achelous/upstream-core/config"

	"github.com/oklog/ulid"
)

func upload(cdata *config.Config, queue commonQueue.QueueRef, id ulid.ULID) error {
	log.Printf("fake-uploading message %s of queue /%s", id, queue)
	return nil
}
