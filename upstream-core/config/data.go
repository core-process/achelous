package config

import (
	"github.com/vrischmann/jsonutil"
)

type Config struct {
	PauseBetweenRuns struct {
		PreviousRunOK         jsonutil.Duration
		PreviousRunWithErrors jsonutil.Duration
	}
	Target struct {
		Upload struct {
			URL    *string
			Header map[string][]string
		}
		Report struct {
			URL    *string
			Header map[string][]string
		}
		RetriesPerRun struct {
			Attempts             int
			PauseBetweenAttempts jsonutil.Duration
		}
	}
}
