package main

import (
	"fmt"
	"os"

	"github.com/core-process/achelous/spring-core/args"
	"github.com/core-process/achelous/spring-core/programs"
)

func main() {
	// parse arguments
	program, smArgs, mqArgs, values, err := args.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// run sub programs
	switch program {
	case args.ArgProgramSendmail:
		err = programs.Sendmail(smArgs, values)
	case args.ArgProgramMailq:
		err = programs.Mailq(mqArgs)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	os.Exit(0)
}
