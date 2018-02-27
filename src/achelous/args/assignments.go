package args

var assignements = map[string]interface{}{
	"*string": func(source string, target **string) {
		*target = &source
	},
}
