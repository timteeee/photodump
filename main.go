package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"photodump/internal/app"
	"photodump/internal/env"
	"photodump/internal/options"

	flag "github.com/spf13/pflag"
)

var (
	dev  *bool   = flag.Bool("dev", false, "whether to run in dev mode")
	port *uint16 = flag.Uint16P("port", "p", 80, "port to bind to")
)

type Env struct {
	StorageEndpoint  string `env:"storage_endpoint"`
	StorageAccessKey string `env:"storage_access_key"`
	StorageSecretKey string `env:"storage_secret_key"`
	StorageBucket    string `env:"storage_bucket"`
	DBUser           string `env:"db_user"`
	DBPassword       string `env:"db_password"`
	DBHost           string `env:"db_host"`
	DBPort           uint16 `env:"db_port"`
	DBDatabase       string `env:"db_database"`
}

var (
	//go:embed templates
	templates embed.FS

	//go:embed static
	static embed.FS

	//go:embed migrations
	migrations embed.FS
)

func main() {
	flag.Parse()

	opts := options.Default()
	opts.Port = *port

	if *dev {
		opts.Mode = options.DEV
		options.LogLevel.Set(slog.LevelDebug)
	}

	e := &Env{}
	env.MustParse(e)

	opts.StorageEndpoint = e.StorageEndpoint
	opts.StorageAccessKey = e.StorageAccessKey
	opts.StorageSecretKey = e.StorageSecretKey
	opts.Bucket = e.StorageBucket

	opts.DBUser = e.DBUser
	opts.DBPassword = e.DBPassword
	opts.DBHost = e.DBHost
	opts.DBPort = e.DBPort
	opts.DBDatabase = e.DBDatabase

	runCtx, runCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer runCancel()

	if err := app.New(opts, templates, migrations, static).Run(runCtx); err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}
}
