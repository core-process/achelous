package processor

import "log"

func report(OK bool) error {
	log.Printf("fake-report OK=%v", OK)
	return nil
}
