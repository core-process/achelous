package main

import (
	"io"
	"log"
	"log/syslog"
	"os"

	"github.com/core-process/achelous/spring-core/args"
	"github.com/core-process/achelous/spring-core/config"
	"github.com/core-process/achelous/spring-core/debug"
	"github.com/core-process/achelous/spring-core/programs"
)

func main() {
	// configure logger to write to the syslog
	logwriter, err := syslog.New(syslog.LOG_DEBUG|syslog.LOG_MAIL, "achelous/spring-core")
	if err != nil {
		log.Println(err)
		return
	}
	log.SetOutput(io.MultiWriter(logwriter, os.Stderr))

	// load config
	cdata, err := config.Load()
	if err != nil {
		log.Printf("failed to load spring config: %v", err)
	}

	// parse arguments
	program, smArgs, mqArgs, values, err := args.Parse(os.Args)
	if err != nil {
		log.Printf("failed to parse arguments: %v", err)
		debug.MailOn(debug.InvalidParameters, cdata, "parsing arguments", nil, err, os.Args)
		os.Exit(1)
	}

	// run sub programs
	if program == args.ArgProgramSendmail && smArgs.Arg_bp {
		mqArgs = new(args.MqArgs)
		mqArgs.Arg_v = smArgs.Arg_v
		err = programs.Mailq(mqArgs)
	} else {
		switch program {
		case args.ArgProgramSendmail:
			err = programs.Sendmail(cdata, smArgs, values)
		case args.ArgProgramMailq:
			err = programs.Mailq(mqArgs)
		}
	}

	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	os.Exit(0)
}
