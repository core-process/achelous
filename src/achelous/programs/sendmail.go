package programs

import (
	"achelous/args"
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/jhillyerd/enmime"
)

func Sendmail(smArgs *args.SmArgs, recipients []string) {
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

	// generate a uuid for message
	id, err := uuid.NewUUID()
	if id == uuid.Nil || err != nil {
		fmt.Print(err)
		return
	}

	// parse envelope
	envelope, err := enmime.ReadEnvelope(stdin)
	if err != nil {
		fmt.Print(err)
		return
	}

	addresses, _ := envelope.AddressList("From")
	for _, address := range addresses {
		fmt.Printf("From: %s <%s>\n", address.Name, address.Address)
	}

	addresses, _ = envelope.AddressList("To")
	for _, address := range addresses {
		fmt.Printf("To: %s <%s>\n", address.Name, address.Address)
	}

	fmt.Printf("Subject: %q\n", envelope.GetHeader("Subject"))
	fmt.Println()

	fmt.Printf("Text Body: %q\n", envelope.Text)
	fmt.Printf("HTML Body: %q\n", envelope.HTML)
	fmt.Println()

	fmt.Println("Envelope errors:")
	for _, e := range envelope.Errors {
		fmt.Println(e.String())
	}
}
