package main

import (
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"syscall"
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
	signal_chan := make(chan os.Signal, 1)
	signal.Notify(signal_chan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// create exit channel and listen to signal channel
	exit_chan := make(chan int)
	go func() {
		for {
			s := <-signal_chan
			switch s {
			// kill -SIGHUP XXXX
			case syscall.SIGHUP:
				log.Println("SIGHUP")

			// kill -SIGINT XXXX or Ctrl+c
			case syscall.SIGINT:
				log.Println("SIGINT")
				exit_chan <- 0

			// kill -SIGTERM XXXX
			case syscall.SIGTERM:
				log.Println("SIGTERM")
				exit_chan <- 0

			// kill -SIGQUIT XXXX
			case syscall.SIGQUIT:
				log.Println("SIGQUIT")
				exit_chan <- 0

			default:
				log.Println("UNKNOWN")
				exit_chan <- 1
			}
		}
	}()

	// listen to exit channel
	code := <-exit_chan
	os.Exit(code)
}
