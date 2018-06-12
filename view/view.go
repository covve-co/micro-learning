package view

import (
	"errors"
	"html/template"
	"io"

	"gitlab.com/covveco/special-needs/view/content"
	"gitlab.com/covveco/special-needs/view/layouts"
)

var templates map[string]*template.Template

func Parse() error {
	templates = make(map[string]*template.Template)

	for _, name := range content.AssetNames() {
		b, err := content.Asset(name)
		if err != nil {
			return err
		}

		t := template.New(name).Funcs(funcs)

		t, err = t.Parse(string(b))
		if err != nil {
			return err
		}

		for _, name := range layouts.AssetNames() {
			b, err := layouts.Asset(name)
			if err != nil {
				return err
			}

			t, err = t.Parse(string(b))
			if err != nil {
				return err
			}
		}

		templates[name] = t
	}

	return nil
}

func Render(w io.Writer, name string, data interface{}) error {
	t, ok := templates[name+".html"]
	if !ok {
		return errors.New("no template found")
	}

	return t.Execute(w, data)
}
