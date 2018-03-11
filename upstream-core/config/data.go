package config

import (
	"net/url"
	"time"
)

type Config struct {
	PauseBetweenRuns time.Duration
	Target           struct {
		Upload struct {
			URL    *url.URL
			Query  map[string]string
			Header map[string]string
		}
		Report struct {
			URL    *url.URL
			Query  map[string]string
			Header map[string]string
		}
		RetryPerRun struct {
			Attempts int
			Pause    time.Duration
		}
	}
}
