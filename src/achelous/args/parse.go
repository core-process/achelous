package args

import (
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
	field reflect.StructField
	value reflect.Value
}

func Parse(argv []string) Args {
	// generate configuration
	result := Args{}

	argsType := reflect.TypeOf(result)
	argsValue := reflect.ValueOf(&result).Elem()

	config := make([]argConf, argsValue.NumField())

	for i := 0; i < argsType.NumField(); i++ {
		config[i].field = argsType.Field(i)
		config[i].value = argsValue.FieldByName(config[i].field.Name)
		config[i].name = "-" + config[i].field.Name[len("Arg_"):]
		switch config[i].field.Tag.Get("args") {
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

				lookup := config[ci].field.Type.Name()
				if config[ci].field.Type.Kind() == reflect.Ptr {
					lookup = "*" + config[ci].field.Type.Elem().Name()
				}

				reflect.
					ValueOf(assignements[lookup]).
					Call([]reflect.Value{
						reflect.ValueOf(source),
						config[ci].value.Addr(),
					})
			}

			switch config[ci]._type {
			case argTypeFlag:
				if config[ci].name == argv[ai] {
					config[ci].value.SetBool(true)
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
