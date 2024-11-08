package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	flag "github.com/spf13/pflag"
)

var (
	port  *uint16 = flag.Uint16P("port", "p", 80, "port to bind to")
	title *string = flag.StringP("title", "t", "PhotoDump", "title for the application")
)

type Layout struct {
	Title string
}

func main() {
	flag.Parse()

	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *port),
		Handler: router(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		fmt.Print("\r")
		slog.Info("Gracefully shutting down...")

		deadline := time.Second * 30
		shutdownCtx, _ := context.WithTimeout(ctx, deadline)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("Graceful shutdown timed out. Forcing exit.")
			}

			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Fatalf("error while shutting down server %s", err.Error())
			}
		}()
		cancel()
	}()

	slog.Info(fmt.Sprintf("Listening on port %d...", *port))
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("uh oh %s", err.Error())
	}

	slog.Info("Server shut down successfully.")

	<-ctx.Done()
}

func router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		layout := filepath.Join("templates", "layout.html")
		home := filepath.Join("templates", "index.html")

		t, err := template.ParseFiles(layout, home)
		if err != nil {
			http.Error(w, fmt.Sprintf("I goofed: %s", err.Error()), 500)
		}
		l := Layout{Title: *title}
		t.ExecuteTemplate(w, "layout", l)
	})

	return r
}
