package queue

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"
	"syscall"

	"github.com/core-process/achelous/common/config"

	ps "github.com/mitchellh/go-ps"
)

func Trigger() error {

	// read pid file
	pidraw, err := ioutil.ReadFile(config.UpstreamPid)
	if err != nil {
		return err
	}

	pid, err := strconv.ParseUint(string(pidraw), 10, 64)
	if err != nil {
		return err
	}

	// send SIGUSR1 signal
	err = syscall.Kill(int(pid), syscall.SIGUSR1)
	if err != nil {
		log.Printf("failed to send SIGUSR1 to primary process: %v", err)

		// retrieve list of all processes
		log.Printf("trying to send SIGUSR1 to child processes...")
		procs, err := ps.Processes()
		if err != nil {
			return err
		}

		// try to send signal to child processes
		sentToChild := false

		for _, proc := range procs {
			// compare parent pid to primary process pid
			if proc.PPid() == int(pid) {
				// send SIGUSR1 to child process
				err = syscall.Kill(proc.Pid(), syscall.SIGUSR1)
				if err != nil {
					log.Printf("failed to send SIGUSR1 to child process %d: %v", proc.Pid(), err)
				} else {
					sentToChild = true
				}
			}
		}

		if !sentToChild {
			return errors.New("Failed to send SIGUSR1 to any process")
		}
	}

	return nil
}
