package programs

import (
	"achelous/args"
	"bufio"
	"fmt"
	"io"
	"os"

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

	// parse envelope
	env, err := enmime.ReadEnvelope(stdin)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("From: %v\n", env.GetHeader("From"))

	alist, _ := env.AddressList("To")
	for _, addr := range alist {
		fmt.Printf("To: %s <%s>\n", addr.Name, addr.Address)
	}

	fmt.Printf("Root Part Subject: %q\n", env.Root.Header.Get("Subject"))
	fmt.Printf("Envelope Subject: %q\n", env.GetHeader("Subject"))
	fmt.Println()

	fmt.Printf("Text Content: %q\n", env.Text)
	fmt.Printf("HTML Content: %q\n", env.HTML)
	fmt.Println()

	fmt.Println("Envelope errors:")
	for _, e := range env.Errors {
		fmt.Println(e.String())
	}
}
