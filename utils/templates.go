package utils

import (
	"bytes"
	"text/template"
)

type cacheMap map[string]*template.Template

var stringsCache = make(cacheMap)
var filesCache = make(cacheMap)

func process(t *template.Template, vars interface{}) (string, error) {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		return "", err
	}

	return tmplBytes.String(), nil
}

func ProcessTemplateString(name string, str string, vars interface{}) (string, error) {
	var err error

	tmpl := stringsCache[name]

	if tmpl != nil {
		err = nil
	} else {
		tmpl, err = template.New(name).Parse(str)

		if err != nil {
			return "", err
		}
	}

	return process(tmpl, vars)
}

func ProcessTemplateFile(fileName string, vars interface{}) (string, error) {
	var err error

	tmpl := filesCache[fileName]

	if tmpl != nil {
		err = nil
	} else {
		tmpl, err = template.ParseFiles(fileName)

		if err != nil {
			return "", err
		}
	}

	filesCache[fileName] = tmpl

	return process(tmpl, vars)
}
