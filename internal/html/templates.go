package html

import (
	"embed"
	"html/template"
	"io"
	"time"
)

//go:embed *.html
var templates embed.FS

var tmpl *template.Template = template.Must(template.New("tmpl").Funcs(funcs).ParseFS(templates, "*.html"))

type Renderer interface {
	Render(io.Writer) error
}

type HomePage struct{}

func (h *HomePage) Render(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "index.html", h)
}

type PhotoCard struct {
	Src        string
	Filename   string
	UploadedAt time.Time
}

func (p *PhotoCard) Render(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "photo-card.html", p)
}
