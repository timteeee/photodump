package html

import (
	"embed"
	"html/template"
	"io"
	"time"
)

//go:embed all:*.html
var templates embed.FS

// Templates
var (
	// Pages
	homeTemplate *template.Template = parseWithLayout("index.html")

	// HTMX Partials
	photoCardTemplate *template.Template = parse("photo-card.html")
)

func parse(filename string) *template.Template {
	return template.Must(template.New(filename).Funcs(funcs).ParseFS(templates, filename))
}

func parseWithLayout(filename string) *template.Template {
	return template.Must(template.New(filename).Funcs(funcs).ParseFS(templates, filename, "layout.html"))
}

type Renderer interface {
	Render(io.Writer) error
}

type HomePage struct{}

func (h *HomePage) Render(w io.Writer) error {
	return homeTemplate.ExecuteTemplate(w, "index.html", h)
}

type PhotoCard struct {
	Src        string
	Filename   string
	UploadedAt time.Time
}

func (pc *PhotoCard) Render(w io.Writer) error {
	return photoCardTemplate.Execute(w, pc)
}
