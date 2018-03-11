package processor

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/core-process/achelous/upstream-core/config"
)

func report(cdata *config.Config, OK bool) error {

	// currently we do report on OK case only
	if !OK {
		return nil
	}

	// check if we have an url available
	if cdata.Target.Report.URL == nil {
		log.Printf("no report target url configured, skipping")
		return nil
	}

	// retry loop
	var lastError error
	lastError = nil

	for i := 0; i < cdata.Target.RetryPerRun.Attempts; i++ {

		if i > 0 {
			time.Sleep(cdata.Target.RetryPerRun.Pause)
			log.Printf("retrying to report: %d", i)
		}

		// prepare request
		req, _ := http.NewRequest("POST", *cdata.Target.Report.URL, nil)
		req.Header = cdata.Target.Report.Header

		// perform request
		res, err := client.Do(req)
		if err != nil {
			// retry
			lastError = err
			continue
		}

		defer res.Body.Close()

		// check status code
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			lastError = errors.New("Invalid status code")
			continue
		}

		// done
		lastError = nil
		break
	}

	return lastError
}
