package processor

import (
	"container/list"
	"context"
	"log"
	"os"
	"path"
	"strings"

	"github.com/core-process/achelous/common/config"
	queuepkg "github.com/core-process/achelous/common/queue"
)

func Run(ctx context.Context) {

	// only true if everything went fine
	allOk := true

	// list of queues to be processes
	queues := list.New()
	queues.PushBack("")

	// upload channel
	cupload := make(chan [2]string)

	go func() {
		for item := range cupload {
			log.Printf("fake-uploading message %s of queue %s", item[1], item[0])
		}
	}()

	// read file queues
	pext := "." + string(queuepkg.MessageStatusPreparing)
	qext := "." + string(queuepkg.MessageStatusQueued)

	for queues.Len() > 0 {
		// check if we have to exit early
		select {
		case <-ctx.Done():
			log.Printf("cancelling current queue walk")
			allOk = false
			break
		default:
			// noop <=> non-blocking
		}

		// pop first element
		queue := queues.Remove(queues.Front()).(string)

		// open directory
		dir, err := os.Open(path.Join(config.Spool, queue))
		if err != nil {
			log.Printf("could not open queue %s: %v", queue, err)
			allOk = false
			continue
		}
		defer dir.Close()

		// read directory
		entries, err := dir.Readdirnames(-1)
		if err != nil {
			log.Printf("could not read entries from queue %s: %v", queue, err)
			allOk = false
			continue
		}

		// iterate entries
		for _, entry := range entries {
			// check if we have to exit early
			select {
			case <-ctx.Done():
				log.Printf("cancelling current queue walk")
				allOk = false
				break
			default:
				// noop <=> non-blocking
			}

			// get file info
			stat, err := os.Stat(path.Join(config.Spool, queue, entry))
			if err != nil {
				log.Printf("could not retrieve file info for entry %s in queue %s: %v", entry, queue, err)
				if !strings.HasSuffix(entry, pext) {
					// in case the error occured while stat'ing a potentially item in preparing
					// state, we will not include this as an invalid operation. this might happen
					// due to race conditions, which are happening by design in this case.
					allOk = false
				}
				continue
			}

			// handle entry
			if stat.Mode().IsDir() {
				// push to list of queues
				queues.PushBack(path.Join(queue, entry))

			} else if stat.Mode().IsRegular() {
				// push to upload channel (if queued item)
				if strings.HasSuffix(entry, qext) {
					id := entry[0 : len(entry)-len(qext)]
					cupload <- [2]string{queue, id}
				}
			}
		}
	}

	// do not report success in case something did not work fine
	if !allOk {
		log.Printf("at least one error occured, skip final report")
	}
}
