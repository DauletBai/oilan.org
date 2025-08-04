// oilan/internal/view/view.go
package view

import (
	"html/template"
	"io"
	"time"
)

// Template represents a single template that can be rendered.
type Template struct {
	htmlTpl *template.Template
}

// NewTemplate creates and parses a set of template files.
func NewTemplate(files ...string) (*Template, error) {
	// We add a custom function 'currentYear' that will be available in all templates.
	funcMap := template.FuncMap{
		"currentYear": func() int {
			return time.Now().Year()
		},
	}

	tpl, err := template.New("base.html").Funcs(funcMap).ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return &Template{htmlTpl: tpl}, nil
}

// Render executes the template with the given data.
func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	return t.htmlTpl.ExecuteTemplate(w, name, data)
}