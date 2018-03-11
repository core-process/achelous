package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	commonConfig "github.com/core-process/achelous/common/config"
)

func Load() (*Config, error) {

	// set default data
	var data Config
	data.DefaultQueue = ""
	data.PrettyJSON = true

	// read file
	raw, err := ioutil.ReadFile(commonConfig.SpringConfig)
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