package config

import (
	"github.com/coreprocess/achelous/common/queue"
)

type Config struct {
	AbortOnParseErrors bool
	DefaultQueue       string
	PrettyJSON         bool
	TriggerQueueRun    bool
	GenerateDebugMail  struct {
		OnInvalidParameters bool
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
