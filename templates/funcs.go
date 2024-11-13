package html

import (
	"fmt"
	"html/template"
)

var funcs template.FuncMap = template.FuncMap{
	"quote": quote,
}

func quote(val any) string {
	return fmt.Sprintf("%q", val)
}
