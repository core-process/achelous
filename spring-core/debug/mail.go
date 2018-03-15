package debug

import (
	"github.com/core-process/achelous/spring-core/config"

	"github.com/oklog/ulid"
)

type ErrorCategory int8

const (
	InvalidParameters ErrorCategory = 1
	ParsingErrors     ErrorCategory = 2
	OtherErrors       ErrorCategory = 3
)

func MailOn(category ErrorCategory, cdata *config.Config, ref *ulid.ULID, err error, data interface{}) {
}
