package programs

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/core-process/achelous/common/queue"
	"github.com/core-process/achelous/spring-core/args"

	"github.com/jhillyerd/enmime"
)

func Sendmail(smArgs *args.SmArgs, recipients []string) error {

	// filter stdin and stop at single-dot-line
	ignoreDot := smArgs.Arg_i || smArgs.Arg_oi || smArgs.Arg_O.Opt_IgnoreDots

	filterReader, filterWriter := io.Pipe()

	go func() {
		// close writer at end of func
		defer filterWriter.Close()

		// scan lines, perform modifications
		scanner := bufio.NewScanner(os.Stdin)
		headerSection := true
		currentLine := -1

		for scanner.Scan() {
			// get next line
			line := scanner.Text()
			currentLine++
			// handle "."
			if !ignoreDot && line == "." {
				break
			}
			// handle end of header
			if headerSection {
				if len(line) == 0 {
					// we leave the header section naturally
					headerSection = false
				} else {
					// we have to simulate the end of the header section
					isHeaderLine := strings.Index(line, ":") > 0
					isContLine := strings.IndexAny(line, " \t") == 0 && currentLine > 0

					if !isHeaderLine && !isContLine {
						line = "\n" + line
						headerSection = false
					}
				}
			}
			// forward data (and errors)
			line += "\n"
			if _, err := io.WriteString(filterWriter, line); err != nil {
				filterWriter.CloseWithError(err)
				break
			}
		}
		// forward errors
		if err := scanner.Err(); err != nil {
			filterWriter.CloseWithError(err)
		}
	}()

	// parse envelope
	envelope, err := enmime.ReadEnvelope(filterReader)
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

	// enhance From and To
	if len(envelope.Root.Header.Get("From")) == 0 {
		if smArgs.Arg_f != nil && smArgs.Arg_F != nil {
			envelope.Root.Header.Set("From", *smArgs.Arg_F+" <"+*smArgs.Arg_f+">")
		} else if smArgs.Arg_f != nil {
			envelope.Root.Header.Set("From", *smArgs.Arg_f)
		} else if smArgs.Arg_F != nil {
			envelope.Root.Header.Set("From", *smArgs.Arg_F)
		}
	}

	if len(envelope.Root.Header.Get("To")) == 0 && len(recipients) > 0 {
		envelope.Root.Header.Set("To", strings.Join(recipients, ", "))
	}

	// add to queue
	err = queue.Add(queue.QueueRefRoot, envelope)
	if err != nil {
		return err
	}

	return nil
}
