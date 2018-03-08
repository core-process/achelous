package queue

import (
	"path"

	"github.com/core-process/achelous/common/config"

	"github.com/oklog/ulid"
)

func QueuePath(queue QueueRef) string {
	return path.Join(config.Spool, string(queue))
}

func MessagePath(queue QueueRef, msgId ulid.ULID, status MessageStatus) string {
	return path.Join(QueuePath(queue), msgId.String()+"."+string(status))
}
