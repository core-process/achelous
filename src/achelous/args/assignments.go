package args

import (
	"errors"
	"strings"
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
			return errors.New("invalid value (" + source + ")")
		}
		*target = &value
		return nil
	},

	"*Arg_N": func(source string, target **Arg_N) error {
		value := Arg_N_Never
		for _, flag := range strings.Split(source, ",") {
			switch flag {
			case "never":
				break // ignore
			case "failure":
				value |= Arg_N_Failure
			case "delay":
				value |= Arg_N_Delay
			case "success":
				value |= Arg_N_Success
			default:
				return errors.New("invalid value (" + flag + ")")
			}
		}
		*target = &value
		return nil
	},

	"*Arg_p": func(source string, target **Arg_p) error {
		value := Arg_p{}
		parts := strings.Split(source, ":")
		switch len(parts) {
		case 2:
			value.Hostname = &parts[1]
			fallthrough
		case 1:
			value.Protocol = parts[0]
		default:
			return errors.New("invalid value (" + source + ")")
		}
		*target = &value
		return nil
	},

	"*Arg_R": func(source string, target **Arg_R) error {
		var value Arg_R
		switch source {
		case "full":
			value = Arg_R_Full
		case "hdrs":
			value = Arg_R_Hdrs
		default:
			return errors.New("invalid value (" + source + ")")
		}
		*target = &value
		return nil
	},
}
