package queue

import (
	"io/ioutil"
	"strconv"
	"syscall"

	"github.com/core-process/achelous/common/config"
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
		return err
	}

	return nil
}
