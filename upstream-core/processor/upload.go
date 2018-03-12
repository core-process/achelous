package processor

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	commonQueue "github.com/core-process/achelous/common/queue"
	"github.com/core-process/achelous/upstream-core/config"

	"github.com/oklog/ulid"
)

func upload(cdata *config.Config, queue commonQueue.QueueRef, id ulid.ULID) error {

	// retry loop
	var lastError error
	lastError = nil

	for i := 0; i < cdata.Target.RetriesPerRun.Attempts; i++ {

		if i > 0 {
			time.Sleep(cdata.Target.RetriesPerRun.PauseBetweenAttempts.Duration)
			log.Printf("retrying to upload message %s in queue /%s (i=%d)", id, queue, i)
		}

		// open message file
		fd, err := os.Open(commonQueue.MessagePath(queue, id, commonQueue.MessageStatusQueued))
		if err != nil {
			return err
		}

		defer fd.Close()

		// prepare request
		req, _ := http.NewRequest("POST", *cdata.Target.Upload.URL, fd)
		req.Header = cdata.Target.Upload.Header
		req.Header.Set("Content-Type", "application/json")

		// perform request
		res, err := client.Do(req)
		if err != nil {
			// retry
			lastError = err
			continue
		}

		defer res.Body.Close()

		// check status code
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			lastError = errors.New("Invalid status code")
			continue
		}

		// done
		lastError = nil
		break
	}

	return lastError
}
