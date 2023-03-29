package html

import (
	"embed"
	"html/template"
)

//go:embed *.html
var question embed.FS

func GetQuestionTemplate() *template.Template {
	return template.Must(template.New("").ParseFS(question, "*"))
}
