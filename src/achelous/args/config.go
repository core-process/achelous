package args

import "reflect"

type argType int8

const (
	argTypeFlag  = 1
	argTypeValue = 2
)

type argConf struct {
	name  string
	_type argType
	value reflect.Value
}

func prepareConfig(argsType reflect.Type, argsValue reflect.Value) []argConf {

	config := make([]argConf, argsValue.NumField())

	for i := 0; i < argsType.NumField(); i++ {
		field := argsType.Field(i)
		config[i].value = argsValue.FieldByName(field.Name)
		config[i].name = "-" + field.Name[len("Arg_"):]
		switch field.Tag.Get("args") {
		case "flag":
			config[i]._type = argTypeFlag
		case "value":
			config[i]._type = argTypeValue
		default:
			panic("Invalid arg type " + field.Tag.Get("args") + " at " + field.Name)
		}
	}

	return config
}
