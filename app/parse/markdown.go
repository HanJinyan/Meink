package parse

import (
	"html/template"

	"github.com/russross/blackfriday/v2"
)

func Markdown(markdown string) template.HTML {
	return template.HTML(blackfriday.Run([]byte(markdown)))
}
