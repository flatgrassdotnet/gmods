package frontend

import (
	"text/template"
)

var t *template.Template

func Init() error {
	var err error

	t, err = template.New("base.html").Funcs(templateFuncs).ParseFiles("templates/base.html")
	if err != nil {
		return err
	}

	t, err = t.ParseGlob("templates/include/*.html")
	if err != nil {
		return err
	}

	return nil
}
