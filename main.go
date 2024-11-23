package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"photodump/internal/app"

	"github.com/go-chi/httplog/v2"
	flag "github.com/spf13/pflag"
)

var (
	dev     *bool   = flag.Bool("dev", false, "whether to run in dev mode")
	port    *uint16 = flag.Uint16P("port", "p", 80, "port to bind to")
	storage *string = flag.StringP("storage", "s", "", "URL for object storage")
	bucket  *string = flag.StringP("bucket", "b", "", "bucket for object storage")
	dbUrl   *string = flag.String("db", "", "URL for database")
)

var (
	//go:embed templates
	templates embed.FS

	//go:embed static
	static embed.FS
)

func main() {
	flag.Parse()

	// TODO: get these from env vars
	accessKey := "VPP0fkoCyBZx8YU0QTjH"
	secretKey := "iFq6k8RLJw5B0faz0cKCXeQk0w9Q8UdtaFzHuw4J"

	logLevel := slog.LevelInfo
	if *dev {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(httplog.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	opts := app.Options{
		Port:                 *port,
		Logger:               logger,
		Templates:            templates,
		StaticAssets:         static,
		ObjectStoreEndpoint:  *storage,
		Bucket:               *bucket,
		Secure:               !*dev,
		ObjectStoreAccessKey: accessKey,
		ObjectStoreSecretKey: secretKey,
	}

	runCtx, runCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer runCancel()

	if err := app.New(&opts).Run(runCtx); err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}
}
