package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Printf("Hello, world.\n")
	for _, e := range os.Environ() {
		fmt.Println(e)
	}
}
