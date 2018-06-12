package view

import "html/template"

// funcs contains all the functions to be used in templates.
var funcs = template.FuncMap{
	"inc": inc,
}

func inc(i int) int {
	return i + 1
}
