package programs

import (
	"achelous/args"
	"fmt"
	"os"

	"github.com/jhillyerd/enmime"
)

// https://github.com/jhillyerd/enmime
// https://github.com/veqryn/go-email
// https://github.com/DusanKasan/parsemail

func Sendmail(smArgs *args.SmArgs, recipients []string) {
	// read mail from stdin
	//ignoreDots := smArgs.Arg_i || smArgs.Arg_O.Opt_IgnoreDots
	/*
		lines := []string{}
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			if !ignoreDots && line == "." {
				break
			}
			lines = append(lines, line)
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading standard input:", err)
			return
		}
	*/

	env, err := enmime.ReadEnvelope(os.Stdin)
	if err != nil {
		fmt.Print(err)
		return
	}

	// Headers can be retrieved via Envelope.GetHeader(name).
	fmt.Printf("From: %v\n", env.GetHeader("From"))

	// Address-type headers can be parsed into a list of decoded mail.Address structs.
	alist, _ := env.AddressList("To")
	for _, addr := range alist {
		fmt.Printf("To: %s <%s>\n", addr.Name, addr.Address)
	}

	// enmime can decode quoted-printable headers.
	fmt.Printf("Subject: %v\n", env.GetHeader("Subject"))

	// The plain text body is available as mime.Text.
	fmt.Printf("Text Body: %v chars\n", len(env.Text))

	// The HTML body is stored in mime.HTML.
	fmt.Printf("HTML Body: %v chars\n", len(env.HTML))

	// mime.Inlines is a slice of inlined attacments.
	fmt.Printf("Inlines: %v\n", len(env.Inlines))

	// mime.Attachments contains the non-inline attachments.
	fmt.Printf("Attachments: %v\n", len(env.Attachments))
}
