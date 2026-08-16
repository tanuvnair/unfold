// Package webui embeds the production Vite build and serves it as a SPA.
//
// go:embed cannot reach outside this package, so Makefile copies web/dist
// into dist/ before building the API binary. When only .gitkeep is present,
// HasUI reports false and the API stays API-only.
package webui

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed all:dist
var dist embed.FS

// HasUI reports whether a production index.html was embedded.
func HasUI() bool {
	f, err := dist.Open("dist/index.html")
	if err != nil {
		return false
	}
	_ = f.Close()
	return true
}

// Handler serves embedded static files. Unknown paths fall back to
// index.html so TanStack Router client routes keep working.
func Handler() http.Handler {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		return http.NotFoundHandler()
	}
	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		f, err := sub.Open(path)
		if err != nil {
			if !HasUI() {
				http.NotFound(w, r)
				return
			}
			r = r.Clone(r.Context())
			r.URL.Path = "/"
			fileServer.ServeHTTP(w, r)
			return
		}
		_ = f.Close()
		fileServer.ServeHTTP(w, r)
	})
}
