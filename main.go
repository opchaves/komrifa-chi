package main

import (
	"context"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jritsema/gotoolbox/web"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed all:public/*
	publicFS embed.FS

	//parsed templates
	html *template.Template
)

func main() {
	// The HTTP Server
	server := &http.Server{Addr: ":8080", Handler: service()}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func service() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	var err error
	html, err = web.TemplateParseFSRecursive(templateFS, ".html", true, nil)
	if err != nil {
		panic(err)
	}

	r.Get("/", index)
	r.Get("/company", companies)
	r.Post("/company", createCompany)
	r.Get("/company/add", companyAdd)
	r.Get("/company/{id}", getCompany)
	r.Put("/company/{id}", saveCompany)
	r.Delete("/company/{id}", removeCompany)
	r.Get("/company/edit/{id}", companyEdit)

	// Create a route along /files that will serve contents from the ./data/ folder.
	// workDir, _ := os.Getwd()
	// filesDir := http.Dir(filepath.Join(workDir, "public"))
	// fileServer(r, "/static", filesDir)

	staticContent, err := fs.Sub(fs.FS(publicFS), "public")
	if err != nil {
		log.Fatal(err)
	}
	fileServer(r, "/static", http.FS(staticContent))

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
