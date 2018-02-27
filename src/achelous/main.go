package main

import (
	"achelous/args"
	"fmt"
)

func main() {
	var a = args.Parse(make([]string, 0))
	fmt.Printf("%+v\n", a)
}
