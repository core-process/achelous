package main

import (
	"achelous/args"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	spew.Dump(os.Args)
	var a = args.Parse(os.Args)
	spew.Dump(a)
}
