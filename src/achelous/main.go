package main

import (
	"achelous/args"
	"achelous/programs"
	"fmt"
	"os"
)

func main() {
	// parse arguments
	program, smArgs, mqArgs, values, err := args.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	// run sub programs
	switch program {
	case args.ArgProgramSendmail:
		programs.Sendmail(smArgs, values)
	case args.ArgProgramMailq:
		programs.Mailq(mqArgs)
	}
}
