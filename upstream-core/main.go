package main

import (
	"context"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/core-process/achelous/upstream-core/config"
	"github.com/core-process/achelous/upstream-core/processor"
)

func main() {
	// configure logger to write to the syslog
	logwriter, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_MAIL, "achelous/upstream-core")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(logwriter)

	// create signal channel
	csig := make(chan os.Signal, 1)
	signal.Notify(
		csig,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-csig
		log.Println("exiting... (" + sig.String() + ")")
		cancel()
	}()

	// main loop
	cancelled := false

	for !cancelled {
		// load config (reload on every run to ensure most up-to-date config)
		cdata, err := config.Load()
		if err != nil {
			log.Printf("failed to load upstream config: %v", err)
		}

		// run processor
		runOK := processor.Run(cdata, ctx)

		// select on timeout and signal
		pauseBetweenRuns := cdata.PauseBetweenRuns.PreviousRunWithErrors.Duration
		if runOK {
			pauseBetweenRuns = cdata.PauseBetweenRuns.PreviousRunOK.Duration
		}

		select {
		case <-ctx.Done():
			cancelled = true
		case <-time.After(pauseBetweenRuns):
			// noop
		}
	}

	// completed
	log.Println("core completed")
	os.Exit(0)
}
