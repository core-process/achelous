package args

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var assignements = map[string]interface{}{

	"*string": func(source string, target **string) error {
		*target = &source
		return nil
	},

	"bool": func(source string, target *bool) error {
		var value bool
		switch source {
		case "true", "yes", "1", "":
			value = true
		case "false", "no", "0":
			value = false
		}
		*target = value
		return nil
	},

	"*int32": func(source string, target **rune) error {
		value, _ := utf8.DecodeRuneInString(source)
		if value == utf8.RuneError {
			return errors.New("invalid value (" + source + ")")
		}
		*target = &value
		return nil
	},

	"*int16": func(source string, target **int16) error {
		value64, err := strconv.ParseInt(source, 10, 16)
		if err != nil {
			return err
		}
		value := int16(value64)
		*target = &value
		return nil
	},

	"*Duration": func(source string, target **time.Duration) error {
		var value time.Duration
		num := ""
		unit := 's'
		apply := func() error {
			if num == "" {
				return nil
			}
			numVal, err := strconv.ParseInt(num, 10, 16)
			if err != nil {
				return err
			}
			switch unit {
			case 's':
				value += time.Second * time.Duration(numVal)
			case 'm':
				value += time.Minute * time.Duration(numVal)
			case 'h':
				value += time.Hour * time.Duration(numVal)
			case 'd':
				value += time.Hour * 24 * time.Duration(numVal)
			case 'w':
				value += time.Hour * 24 * 7 * time.Duration(numVal)
			default:
				return errors.New("invalid unit (" + string(unit) + ")")
			}
			num = ""
			unit = 's'
			return nil
		}
		for _, digit := range source {
			if unicode.IsDigit(digit) {
				num += string(digit)
			} else {
				unit = digit
				err := apply()
				if err != nil {
					return err
				}
			}
		}
		err := apply()
		if err != nil {
			return err
		}
		*target = &value
		return nil
	},

	"*Arg_main": func(source string, target **Arg_main) error {
		var value Arg_main
		switch source {
		case "newaliases":
			value = Arg_main_newaliases
		case "mailq":
			value = Arg_main_mailq
		default:
			value = Arg_main_sendmail
		}
		*target = &value
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
		parts := strings.SplitN(source, ":", 2)
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

func init() {
	assignements["Arg_O"] =
		func(source string, target *Arg_O) error {
			// extract data
			parts := strings.SplitN(source, "=", 2)
			name := parts[0]
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			}
			// assign value
			field := reflect.
				ValueOf(target).Elem().
				FieldByName("Opt_" + name)
			return assign(value, field)
		}
}

func assign(source string, target reflect.Value) error {

	_type := target.Type()

	lookup := _type.Name()
	if _type.Kind() == reflect.Ptr {
		lookup = "*" + _type.Elem().Name()
	}

	err, _ :=
		reflect.
			ValueOf(assignements[lookup]).
			Call([]reflect.Value{
				reflect.ValueOf(source),
				target.Addr(),
			})[0].
			Interface().(error)

	return err
}
