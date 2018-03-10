package processor

import (
	"container/list"
	"context"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/core-process/achelous/common/config"
	commonQueue "github.com/core-process/achelous/common/queue"

	"github.com/oklog/ulid"
)

func Run(ctx context.Context) {

	log.Printf("queue processing started")

	// only true if everything went fine
	allOk := true

	// list of queues to be processes
	queues := list.New()
	queues.PushBack("")

	// initiate upload
	jobs := make(chan [2]string)
	var wg sync.WaitGroup

	go func() {
		// manage wait group
		wg.Add(1)
		defer wg.Done()

		// process jobs
		for job := range jobs {
			queueRef := commonQueue.QueueRef(job[0])
			msgID := ulid.MustParse(job[1])
			// run upload
			err := upload(queueRef, msgID)
			// handle errors
			if err != nil {
				log.Printf("upload of message %s in queue /%s failed: %v", job[1], job[0], err)
				allOk = false
				continue
			}
			// remove queue entry
			log.Printf("upload of message %s in queue /%s succeeded", job[1], job[0])
			err = commonQueue.Remove(queueRef, msgID)
			if err != nil {
				log.Printf("could not remove message %s in queue /%s: %v", job[1], job[0], err)
				allOk = false
			}
		}
	}()

	// read file queues
	pext := "." + string(commonQueue.MessageStatusPreparing)
	qext := "." + string(commonQueue.MessageStatusQueued)

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
			log.Printf("could not open queue /%s: %v", queue, err)
			allOk = false
			continue
		}
		defer dir.Close()

		// read directory
		entries, err := dir.Readdirnames(-1)
		if err != nil {
			log.Printf("could not read entries from queue /%s: %v", queue, err)
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
				log.Printf("could not retrieve file info for entry %s in queue /%s: %v", entry, queue, err)
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
					jobs <- [2]string{queue, id}
				}
			}
		}
	}

	// wait for completion of uploads
	close(jobs)
	wg.Wait()

	// do not report success in case something did not work fine
	err := report(allOk)
	if err != nil {
		log.Printf("could not report status: %v", err)
	}

	log.Printf("queue processing completed (allOk=%v)", allOk)
}
