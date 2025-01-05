package templating

import (
	"bytes"
	"html/template"
)

func ExecuteTemplate(str string, data any) (string, error) {
	tmpl, err := template.New("template").Parse(str)
	if err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}
	if err := tmpl.ExecuteTemplate(buf, "template", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
