package processor

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/core-process/achelous/common/config"
)

func Run(ctx context.Context) {
	err := filepath.Walk(config.Spool, func(path string, info os.FileInfo, err error) error {
		// check if we have to exit early
		select {
		case <-ctx.Done():
			log.Printf("Cancelling current queue walk\n")
			return io.EOF
		default:
			// noop <=> non-blocking
		}
		// handle error
		if err != nil {
			log.Printf("prevent panic by handling failure accessing a path %q: %v\n", config.Spool, err)
			return err
		}
		// do some work
		log.Printf("visited file: %q\n", path)
		return nil
	})

	// we had some issues
	if err != nil {
		log.Printf("error walking the path %q: %v\n", config.Spool, err)
	}
}
