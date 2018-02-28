package args

import (
	"fmt"
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

func Parse(argv []string) Args {
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

	// parse arguments
	for ai := 0; ai < len(argv); ai++ {
		// handle arg0
		if ai == 0 {
			for ci := 0; ci < len(config); ci++ {
				if config[ci]._type == argTypeProgram {
					source := filepath.Base(argv[ai])
					config[ci].value.Set(reflect.ValueOf(&source))
					break
				}
			}
			continue
		}
		// handle regular args
		for ci := 0; ci < len(config); ci++ {

			assignValue := func(source string) {
				err := assign(source, config[ci].value)
				if err != nil {
					fmt.Println("Error while parsing " + config[ci].name + ": " + err.Error())
				}
			}

			switch config[ci]._type {
			case argTypeFlag:
				if config[ci].name == argv[ai] {
					assignValue("")
					break
				}
			case argTypeTrailingValue:
				if config[ci].name == argv[ai] {
					assignValue(argv[ai+1])
					ai++
					break
				}
			case argTypeAttachedValue:
				if strings.HasPrefix(argv[ai], config[ci].name) {
					assignValue(argv[ai][len(config[ci].name):])
					break
				}
			}
		}
	}

	return result
}
