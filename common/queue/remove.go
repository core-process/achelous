package queue

import (
	"os"

	"github.com/oklog/ulid"
)

func Remove(queue QueueRef, id ulid.ULID) error {
	ppath := MessagePath(queue, id, MessageStatusPreparing)
	qpath := MessagePath(queue, id, MessageStatusQueued)

	pexists := true
	qexists := true

	_, err := os.Stat(ppath)
	if err != nil {
		if os.IsNotExist(err) {
			pexists = false
		} else {
			return err
		}
	}

	_, err = os.Stat(qpath)
	if err != nil {
		if os.IsNotExist(err) {
			qexists = false
		} else {
			return err
		}
	}

	if pexists == false && qexists == false {
		return os.ErrNotExist
	}

	if pexists {
		err = os.Remove(ppath)
		if err != nil {
			return err
		}
	}

	if qexists {
		err = os.Remove(qpath)
		if err != nil {
			return err
		}
	}

	return nil
}
