package args

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
)

func Parse(argv []string) (ArgProgram, *SmArgs, *MqArgs, []string, error) {

	// prepare result variables
	resultProgram := ArgProgramSendmail
	var resultSmArgs *SmArgs
	var resultMqArgs *MqArgs
	var resultValues []string

	// detect program
	switch filepath.Base(argv[0]) {
	case "newaliases":
		resultProgram = ArgProgramNewaliases
	case "mailq":
		resultProgram = ArgProgramMailq
	}

	// prepare configuration
	config := []argConf{}

	switch resultProgram {
	case ArgProgramSendmail:
		{

			resultSmArgs = new(SmArgs)
			config = prepareConfig(
				reflect.TypeOf(*resultSmArgs),
				reflect.ValueOf(resultSmArgs).Elem())
		}
	case ArgProgramMailq:
		{
			resultMqArgs = new(MqArgs)
			config = prepareConfig(
				reflect.TypeOf(*resultMqArgs),
				reflect.ValueOf(resultMqArgs).Elem())
		}
	}

	// parse arguments
	assignValue := func(ci int, source string) error {
		err := convert(source, config[ci].value)
		if err != nil {
			return errors.New("Error while parsing " + config[ci].name + ": " + err.Error())
		}
		return nil
	}

	for ai := 1; ai < len(argv); ai++ {
		if strings.HasPrefix(argv[ai], "-") {
			// ignore empty arg
			if argv[ai] == "-" || argv[ai] == "--" {
				continue
			}
			// make sure params and free values are not mixed
			if len(resultValues) > 0 {
				return resultProgram,
					resultSmArgs,
					resultMqArgs,
					resultValues,
					errors.New("Param following free values (" + argv[ai] + ")")
			}
			// handle params
			handled := false
			for ci := 0; ci < len(config); ci++ {
				switch config[ci]._type {
				case argTypeFlag:
					if config[ci].name == argv[ai] {
						source := ""
						err := assignValue(ci, source)
						if err != nil {
							return resultProgram,
								resultSmArgs,
								resultMqArgs,
								resultValues,
								err
						}
						handled = true
						break
					}
				case argTypeValue:
					if strings.HasPrefix(argv[ai], config[ci].name) {
						source := argv[ai][len(config[ci].name):]
						if len(source) == 0 && ai+1 < len(argv) && !strings.HasPrefix(argv[ai+1], "-") {
							source = argv[ai+1]
							ai++
						}
						err := assignValue(ci, source)
						if err != nil {
							return resultProgram,
								resultSmArgs,
								resultMqArgs,
								resultValues,
								err
						}
						handled = true
						break
					}
				}
			}
			// verify if argument had been processed
			if !handled {
				return resultProgram,
					resultSmArgs,
					resultMqArgs,
					resultValues,
					errors.New("Unknown argument " + argv[ai])
			}
		} else {
			// handle values
			resultValues = append(resultValues, argv[ai])
		}
	}

	// done
	return resultProgram,
		resultSmArgs,
		resultMqArgs,
		resultValues,
		nil
}
