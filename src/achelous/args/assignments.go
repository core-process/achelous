package args

import (
	"errors"
)

var assignements = map[string]interface{}{

	"*string": func(source string, target **string) error {
		*target = &source
		return nil
	},

	"*Arg_B": func(source string, target **Arg_B) error {
		var value Arg_B
		switch source {
		case "7BIT":
			value = Arg_B_7Bit
		case "8BITMIME":
			value = Arg_B_8BitMime
		default:
			return errors.New("invalid value")
		}
		*target = &value
		return nil
	},
}
