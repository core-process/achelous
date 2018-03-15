package config

import (
	"github.com/core-process/achelous/common/queue"
)

type Config struct {
	AbortOnParseErrors bool
	DefaultQueue       string
	PrettyJSON         bool
	TriggerQueueRun    bool
	GenerateDebugMail  struct {
		OnUnknownParameters bool
		OnParsingErrors     bool
		OnOtherErrors       bool
		Message             struct {
			Sender   queue.Participant
			Receiver queue.Participant
			Subject  string
			Body     string
		}
	}
}
