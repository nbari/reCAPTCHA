package main

import (
	"html/template"
	"strings"
)

var (
	templateMap = template.FuncMap{
		"Lower": func(s string) string {
			return strings.ToLower(s)
		},
	}
)
