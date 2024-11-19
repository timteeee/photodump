package assets

import (
	"embed"
	"net/http"
)

//go:embed static
var static embed.FS

func Static() http.Handler {
	return http.FileServerFS(static)
}
