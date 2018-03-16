package programs

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	commonConfig "github.com/core-process/achelous/common/config"
	commonQueue "github.com/core-process/achelous/common/queue"
	"github.com/core-process/achelous/spring-core/args"
)

func Mailq(mqArgs *args.MqArgs) error {

	// pretty print queue header
	printQueueHeader := func(queue string) {
		out := "QUEUE: /%s\n\n"
		_, err := fmt.Printf(out, queue)
		if err != nil {
			panic(err)
		}
	}

	// pretty print msg info
	refTime := time.Now()

	printMsg := func(queue string, msg *commonQueue.Message, stat os.FileInfo) {

		fmtParticipant := func(p commonQueue.Participant) string {
			res := []string{}
			if len(p.Name) > 0 {
				res = append(res, p.Name)
			}
			if len(p.Email) > 0 {
				res = append(res, "<"+p.Email+">")
			}
			return strings.Join(res, " ")
		}

		relTs := refTime.Sub(msg.Timestamp)
		relTsStr := ""
		if relTs >= time.Hour {
			relTsStr = strconv.FormatInt(int64(math.Floor(float64(relTs)/float64(time.Hour))), 10) + "h"
		} else if relTs >= time.Minute {
			relTsStr = strconv.FormatInt(int64(math.Floor(float64(relTs)/float64(time.Minute))), 10) + "m"
		} else {
			relTsStr = strconv.FormatInt(int64(math.Floor(float64(relTs)/float64(time.Second))), 10) + "s"
		}

		size := ""
		if stat.Size() >= (1024 * 1024) {
			size = strconv.FormatInt(int64(math.Ceil(float64(stat.Size())/float64(1024*1024))), 10) + "M"
		} else {
			size = strconv.FormatInt(int64(math.Ceil(float64(stat.Size())/float64(1024))), 10) + "K"
		}

		from := ""
		if msg.Participants.From != nil {
			from = fmtParticipant(*msg.Participants.From)
		}

		to := []string{}
		for _, p := range msg.Participants.To {
			pf := fmtParticipant(p)
			if len(pf) > 0 {
				to = append(to, pf)
			}
		}

		out := " %4s %4s %s %s\n" + strings.Repeat(" ", 11) + "%s\n\n"
		_, err := fmt.Printf(
			out,
			relTsStr,
			size,
			msg.ID.String(),
			from,
			strings.Join(to, ", "),
		)
		if err != nil {
			panic(err)
		}
	}

	// only true if everything went fine
	OK := true

	// list of queues to be processes
	queues := list.New()
	queues.PushBack("")

	// read file queues
	qext := "." + string(commonQueue.MessageStatusQueued)

	for queues.Len() > 0 {
		// pop first element
		queue := queues.Remove(queues.Front()).(string)

		// print header
		printQueueHeader(queue)

		// open directory
		dir, err := os.Open(path.Join(commonConfig.Spool, queue))
		if err != nil {
			log.Printf("could not open queue /%s: %v", queue, err)
			OK = false
			continue
		}
		defer dir.Close()

		// read directory
		entries, err := dir.Readdirnames(-1)
		if err != nil {
			log.Printf("could not read entries from queue /%s: %v", queue, err)
			OK = false
			continue
		}

		// iterate entries
		for _, entry := range entries {
			// get file info
			stat, err := os.Stat(path.Join(commonConfig.Spool, queue, entry))
			if err != nil {
				// the item might not exist anymore. this might happen due to
				// race conditions, which are happening by design in this case.
				if !os.IsNotExist(err) {
					log.Printf("could not retrieve file info for entry %s in queue /%s: %v", entry, queue, err)
					OK = false
				}
				continue
			}

			// handle entry
			if stat.Mode().IsDir() {
				// push to list of queues
				queues.PushBack(path.Join(queue, entry))

			} else if stat.Mode().IsRegular() {
				// print only queued entries
				if strings.HasSuffix(entry, qext) {
					// read queued item
					msgRaw, err := ioutil.ReadFile(path.Join(commonConfig.Spool, queue, entry))
					if err != nil {
						// the item might not exist anymore. this might happen due to
						// race conditions, which are happening by design in this case.
						if !os.IsNotExist(err) {
							log.Printf("could not read file for entry %s in queue /%s: %v", entry, queue, err)
							OK = false
						}
						continue
					}
					// decode message
					var msg commonQueue.Message
					err = json.Unmarshal(msgRaw, &msg)
					if err != nil {
						log.Printf("could not parse file for entry %s in queue /%s: %v", entry, queue, err)
						continue
					}
					// print message
					printMsg(queue, &msg, stat)
				}
			}
		}
	}

	// return an error if at least one occured
	if !OK {
		return errors.New("at least one error occured while reading the queue (see error logs for more)")
	}

	return nil
}
