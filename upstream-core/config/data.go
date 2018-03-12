package config

import (
	"github.com/vrischmann/jsonutil"
)

type Config struct {
	PauseBetweenRuns jsonutil.Duration
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
			Pause    jsonutil.Duration
		}
	}
}
