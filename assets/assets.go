package assets

import (
	"embed"
	"net/http"
)

//go:embed public
var public embed.FS

func Public() http.Handler {
	return http.FileServerFS(public)
}
