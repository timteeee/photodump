package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"photodump/assets"
	html "photodump/templates"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	flag "github.com/spf13/pflag"
)

var (
	port *uint16 = flag.Uint16P("port", "p", 80, "port to bind to")
)

// TODO: get rid of this bs
var photoStore sync.Map = sync.Map{}

func main() {
	flag.Parse()

	router := router()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: router,
	}

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		fmt.Println("\rGracefully shutting down...")

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

	fmt.Printf("Listening on port %d...\n", *port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("uh oh %s", err.Error())
	}

	<-ctx.Done()
	fmt.Println("Server shut down successfully")
}

func router() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Mount("/public", assets.Public())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		res := html.HomePage{}
		res.Render(w)
	})

	// TODO: this can go away, photos can be requested with pre-signed urls
	r.Get("/photos/{filename}", func(w http.ResponseWriter, r *http.Request) {
		filename := r.PathValue("filename")
		if filename == "" {
			http.Error(w, "got empty filename", http.StatusBadRequest)
		}

		photo, ok := photoStore.Load(filename)
		if !ok {
			http.NotFound(w, r)
		}

		bytes, ok := photo.([]byte)
		if !ok {
			panic(fmt.Sprintf("expected bytes, got %T", bytes))
		}
		w.Write(bytes)
	})

	// TODO: read from db to get filenames in storage, render photo cards with pre-signed urls
	r.Get("/photos", func(w http.ResponseWriter, r *http.Request) {
		photo := filepath.Join("templates", "photo-card.html")

		for range 3 {
			t, err := template.ParseFiles(photo)
			if err != nil {
				http.Error(w, fmt.Sprintf("I goofed: %s", err.Error()), 500)
			}
			t.Execute(w, nil)
		}
	})

	// TODO: store photos named as hash of contents, store filenames, sizes, datetimes, etc in db, render cards with pre-signed urls
	r.Post("/photos", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, fmt.Sprintf("uh oh, %s", err.Error()), 500)
			return
		}

		if r.MultipartForm == nil {
			http.Error(w, "didn't receive multipart form in request", http.StatusUnprocessableEntity)
			return
		}

		headers := r.MultipartForm.File["photos"]
		if headers == nil {
			http.Error(w, "didn't receive photos in request", http.StatusUnprocessableEntity)
			return
		}

		now := time.Now().UTC()

		for _, h := range headers {
			buf := bytes.NewBuffer(make([]byte, 0, h.Size))
			if photo, err := h.Open(); err == nil {
				io.Copy(buf, photo)
				photo.Close()

				photoStore.Store(h.Filename, buf.Bytes())
				slog.Info("stored photo", "filename", h.Filename)

				card := html.PhotoCard{
					// TODO: create and use presigned url here
					Src:        fmt.Sprintf("/photos/%s", h.Filename),
					Filename:   h.Filename,
					UploadedAt: now,
				}
				card.Render(w)
			}
		}
	})

	return r
}
