package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Options struct {
	// Port to bind to
	Port uint16

	// Object storage configuration
	Bucket               string
	Secure               bool
	ObjectStoreEndpoint  string
	ObjectStoreAccessKey string
	ObjectStoreSecretKey string
}

type App struct {
	// Options related to the HTTP server created in the Run method
	port uint16

	// Object Storage
	objectStore *minio.Client
	bucket      string

	// Base context for application runs
	context context.Context
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

	return &App{
		context:     ctx,
		port:        opts.Port,
		objectStore: client,
		bucket:      opts.Bucket,
	}
}

func (a *App) Run() error {
	runCtx, runCancel := context.WithCancel(a.context)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.port),
		Handler: a.Router(),
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		fmt.Println("\rGracefully shutting down...")

		shutdownCtx, _ := context.WithTimeout(runCtx, time.Second*30)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				fmt.Println("Graceful shutdown timed out. Forcing exit.")
				os.Exit(1)
			}

			if err := server.Shutdown(shutdownCtx); err != nil {
				fmt.Printf("error while shutting down server %s\n", err)
				os.Exit(1)
			}
		}()

		runCancel()
	}()

	fmt.Printf("Listening on port %d...\n", a.port)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("uh oh %s\n", err.Error())
		os.Exit(1)
	}

	<-runCtx.Done()
	fmt.Println("Server shut down successfully")

	return nil
}
