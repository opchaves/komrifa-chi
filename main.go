package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stdrifa/server"
	"syscall"
	"time"
)

var (
	//go:embed all:templates/*
	templateFS embed.FS

	//go:embed all:public/*
	publicFS embed.FS
)

func main() {
	ctx := context.Background()

	server := server.New()

	server.TemplateFS = templateFS
	server.PublicFS = publicFS

	server.ConnectDB(ctx)
	server.ParseTemplates(templateFS)

	// The HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", server.Port),
		Handler: server.Router(),
	}

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
		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

// func router() *chi.Mux {
// 	r.Get("/", index)
// 	r.Route("/company", func(r chi.Router) {
// 		r.Get("/", companies)
// 		r.Post("/", createCompany)
// 		r.Get("/add", companyAdd)
// 		r.Get("/{id}", getCompany)
// 		r.Put("/{id}", saveCompany)
// 		r.Delete("/{id}", removeCompany)
// 		r.Get("/{id}/edit", companyEdit)
// 	})
//
// 	return r
// }
