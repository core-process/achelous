package args

import (
	"errors"
	"path/filepath"
	"reflect"
	"strings"
)

type argType int8

const (
	argTypeProgram       = 0
	argTypeFlag          = 1
	argTypeAttachedValue = 2
	argTypeTrailingValue = 3
)

type argConf struct {
	name  string
	_type argType
	value reflect.Value
}

func Parse(argv []string) (*Args, error) {
	// generate configuration
	result := Args{}

	argsType := reflect.TypeOf(result)
	argsValue := reflect.ValueOf(&result).Elem()

	config := make([]argConf, argsValue.NumField())

	for i := 0; i < argsType.NumField(); i++ {
		field := argsType.Field(i)
		config[i].value = argsValue.FieldByName(field.Name)
		config[i].name = "-" + field.Name[len("Arg_"):]
		switch field.Tag.Get("args") {
		case "program":
			config[i]._type = argTypeProgram
		case "attachedValue":
			config[i]._type = argTypeAttachedValue
		case "trailingValue":
			config[i]._type = argTypeTrailingValue
		default:
			config[i]._type = argTypeFlag
		}
	}

	assignValue := func(ci int, source string) error {
		err := assign(source, config[ci].value)
		if err != nil {
			return errors.New("Error while parsing " + config[ci].name + ": " + err.Error())
		}
		return nil
	}

	// parse arguments
	for ai := 0; ai < len(argv); ai++ {
		// handle arg0
		if ai == 0 {
			for ci := 0; ci < len(config); ci++ {
				if config[ci]._type == argTypeProgram {
					source := filepath.Base(argv[ai])
					err := assignValue(ci, source)
					if err != nil {
						return nil, err
					}
					break
				}
			}
			continue
		}
		// handle regular args
		handled := false
		for ci := 0; ci < len(config); ci++ {
			switch config[ci]._type {
			case argTypeFlag:
				if config[ci].name == argv[ai] {
					source := ""
					err := assignValue(ci, source)
					if err != nil {
						return nil, err
					}
					handled = true
					break
				}
			case argTypeTrailingValue:
				if config[ci].name == argv[ai] {
					if ai+1 >= len(argv) {
						return nil, errors.New("Value missing for argument " + argv[ai])
					}
					source := argv[ai+1]
					err := assignValue(ci, source)
					if err != nil {
						return nil, err
					}
					handled = true
					ai++
					break
				}
			case argTypeAttachedValue:
				if strings.HasPrefix(argv[ai], config[ci].name) {
					source := argv[ai][len(config[ci].name):]
					err := assignValue(ci, source)
					if err != nil {
						return nil, err
					}
					handled = true
					break
				}
			}
		}
		// verify if argument had been processed
		if !handled {
			return nil, errors.New("Unknown argument " + argv[ai])
		}
	}

	return &result, nil
}
