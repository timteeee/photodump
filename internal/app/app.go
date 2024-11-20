package app

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"photodump/internal/html"
	"text/template"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Options struct {
	Logger       *slog.Logger
	Port         uint16
	Templates    fs.FS
	StaticAssets fs.FS

	Bucket               string
	Secure               bool
	ObjectStoreEndpoint  string
	ObjectStoreAccessKey string
	ObjectStoreSecretKey string
}

type App struct {
	logger    *slog.Logger
	port      uint16
	static    fs.FS
	templates *template.Template

	objectStore *minio.Client
	bucket      string
}

func New(opts *Options) *App {
	ctx := context.Background()

	client, err := minio.New(opts.ObjectStoreEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.ObjectStoreAccessKey, opts.ObjectStoreSecretKey, ""),
		Secure: opts.Secure,
	})

	if err != nil {
		fmt.Printf("Error occurred while attempting to connect to object storage at %s:\n\t%s\n", opts.ObjectStoreEndpoint, err)
		os.Exit(1)
	}

	ok, err := client.BucketExists(ctx, opts.Bucket)
	if err != nil {
		fmt.Printf("error checking bucket: %s\n", err)
		os.Exit(1)
	}

	if !ok {
		fmt.Printf("bad bucket setup\n")
		os.Exit(1)
	}

	tmpl := template.Must(template.New("templates").Funcs(html.Funcs).ParseFS(opts.Templates, "templates/*"))

	return &App{
		port:        opts.Port,
		logger:      opts.Logger,
		templates:   tmpl,
		static:      opts.StaticAssets,
		objectStore: client,
		bucket:      opts.Bucket,
	}
}

func (app *App) Run(ctx context.Context) error {
	addr := fmt.Sprintf(":%d", app.port)
	server := &http.Server{
		Addr:    addr,
		Handler: app.router(),
	}

	done := make(chan struct{})
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Error("failed to listen and serve", slog.Any("error", err))
			os.Exit(1)
		}
		close(done)
	}()

	app.logger.Info("Server listening", slog.String("addr", addr))
	select {
	case <-done:
		break
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		server.Shutdown(shutdownCtx)
		cancel()
	}

	return nil
}
