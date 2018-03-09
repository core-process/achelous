package main

import (
	"fmt"
	"time"
)

func main() {
	for true {
		fmt.Println("... running core ...")
		time.Sleep(1000 * time.Millisecond)
	}
}
