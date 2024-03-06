package server

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"stdrifa/app"
	"stdrifa/config"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v2"
)

func (s *Server) Router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(15 * time.Second))
	r.Use(httplog.RequestLogger(s.Logger))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	handlers := app.New(s.DB, s.Logger, s.Templ)

	r.Get("/", handlers.Index)
	r.Route("/company", func(r chi.Router) {
		r.Get("/", handlers.CompanyIndex)
		r.Post("/", handlers.CreateCompany)
		r.Get("/add", handlers.AddCompany)
		r.Get("/{id}", handlers.GetCompany)
		r.Get("/{id}/edit", handlers.EditCompany)
		r.Put("/{id}", handlers.UpdateCompany)
		r.Delete("/{id}", handlers.DeleteCompany)
	})

	if config.IsProduction {
		staticContent, err := fs.Sub(fs.FS(s.PublicFS), "public")
		if err != nil {
			panic(err)
		}
		fileServer(r, "/static", http.FS(staticContent))
	} else {
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "public"))
		fileServer(r, "/static", filesDir)
	}

	return r
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
