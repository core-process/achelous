package config

import (
	"time"
)

type Config struct {
	PauseBetweenRuns time.Duration
	Target           struct {
		Upload struct {
			URL    *string
			Header map[string][]string
		}
		Report struct {
			URL    *string
			Header map[string][]string
		}
		RetryPerRun struct {
			Attempts int
			Pause    time.Duration
		}
	}
}
