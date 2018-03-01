package main

import (
	"achelous/args"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	spew.Dump(os.Args)
	args, err := args.Parse(os.Args)
	if err != nil {
		fmt.Println(err)
		return
	}
	spew.Dump(*args)
}
