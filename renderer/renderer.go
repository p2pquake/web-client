package renderer

import (
	"bytes"
	"html/template"
	"io"
	"os"
	"time"
)

func Render(templateFile string, data interface{}) (string, error) {
	f, err := os.ReadFile("./template/" + templateFile)
	if err != nil {
		return "", err
	}

	t, err := template.New("content").Funcs(template.FuncMap{"date": func() string { return time.Now().Format("01/02 15:04:05") }, "gtag": func() string { return os.Getenv("GTM_CONTAINER_ID") }}).Parse(string(f))
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
