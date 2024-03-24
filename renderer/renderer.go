package renderer

import (
	"bytes"
	"html/template"
	"io"
	"os"
)

func Render(templateFile string, data interface{}) (string, error) {
	f, err := os.ReadFile("./template/" + templateFile)
	if err != nil {
		return "", err
	}

	t, err := template.New("content").Parse(string(f))
	if err != nil {
		return "", err
	}

	t, err = t.ParseGlob("./template/*.html")
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	w := io.Writer(&b)
	err = t.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
