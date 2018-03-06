package queue

import (
	"os/user"
	"path"

	"github.com/oklog/ulid"
)

func QueuePath(queue QueueRef) string {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	return path.Join(user.HomeDir, ".achelous/queues", string(queue))
}

func MessagePath(queue QueueRef, msgId ulid.ULID, status MessageStatus) string {
	return path.Join(QueuePath(queue), msgId.String()+"."+string(status))
}
