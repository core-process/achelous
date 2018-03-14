package config

import (
	"github.com/core-process/achelous/common/queue"
)

type Config struct {
	DefaultQueue      string
	PrettyJSON        bool
	TriggerQueueRun   bool
	GenerateDebugMail struct {
		OnUnknownParameters bool
		OnParsingErrors     bool
		OnEmptyInput        bool
		OnOtherError        bool
		Sender              queue.Participant
		Receiver            queue.Participant
	}
}
