package main

import (
	"achelous/args"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	spew.Dump(os.Args)
	program, smArgs, mqArgs, values, err := args.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	spew.Dump(program)
	if smArgs != nil {
		spew.Dump(*smArgs)
	}
	if mqArgs != nil {
		spew.Dump(*mqArgs)
	}
	spew.Dump(values)
}
