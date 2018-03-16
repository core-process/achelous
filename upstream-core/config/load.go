package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	commonConfig "github.com/core-process/achelous/common/config"

	"github.com/vrischmann/jsonutil"
)

func Load() (*Config, error) {

	// set default data
	var data Config
	data.PauseBetweenRuns.PreviousRunOK = jsonutil.FromDuration(5 * time.Minute)
	data.PauseBetweenRuns.PreviousRunWithErrors = jsonutil.FromDuration(30 * time.Second)
	data.Target.RetriesPerRun.Attempts = 3
	data.Target.RetriesPerRun.PauseBetweenAttempts = jsonutil.FromDuration(5 * time.Second)

	// read file
	raw, err := ioutil.ReadFile(commonConfig.UpstreamConfig)
	if err != nil {
		if os.IsNotExist(err) {
			// do not propagate error
			return &data, nil
		}
		return &data, err
	}

	// unmarshal
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return &data, err
	}

	return &data, nil
}
