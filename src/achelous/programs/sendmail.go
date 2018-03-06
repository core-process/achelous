package programs

import (
	"achelous/args"
	"achelous/queue"
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/jhillyerd/enmime"
)

func Sendmail(smArgs *args.SmArgs, recipients []string) error {
	// prepare standard input
	var stdin io.Reader

	if smArgs.Arg_i || smArgs.Arg_O.Opt_IgnoreDots {
		// use stdin directly
		stdin = os.Stdin

	} else {
		// filter stdin and stop at single-dot-line
		pipeReader, pipeWriter := io.Pipe()

		go func() {
			defer pipeWriter.Close()

			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "." {
					break
				}
				if _, err := io.WriteString(pipeWriter, line+"\n"); err != nil {
					pipeWriter.CloseWithError(err)
					break
				}
			}
			if err := scanner.Err(); err != nil {
				pipeWriter.CloseWithError(err)
			}
		}()

		stdin = pipeReader
	}

	// parse envelope
	envelope, err := enmime.ReadEnvelope(stdin)
	if err != nil {
		return err
	}

	// check for errors
	if len(envelope.Errors) > 0 {
		errMsg := "Parsing failed:"
		for _, v := range envelope.Errors {
			errMsg += "\n- " + v.String()
		}
		return errors.New(errMsg)
	}

	// add to queue
	err = queue.AddToQueue(envelope)
	if err != nil {
		return err
	}

	return nil
}
