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

	data.AbortOnParseErrors = false
	data.DefaultQueue = ""
	data.PrettyJSON = false
	data.TriggerQueueRun = true

	data.GenerateDebugMail.OnUnknownParameters = false
	data.GenerateDebugMail.OnParsingErrors = false
	data.GenerateDebugMail.OnOtherErrors = false
	data.GenerateDebugMail.Message.Sender.Name = "Achelous Spring"
	data.GenerateDebugMail.Message.Receiver.Name = "Devops"
	data.GenerateDebugMail.Message.Subject = "ACHELOUS SPRING DEBUG MESSAGE"
	data.GenerateDebugMail.Message.Body = "Reference: %[1]s\nError: %[2]s\nData: %[3]v"

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
