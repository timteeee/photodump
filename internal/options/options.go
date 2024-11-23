package options

import (
	"log/slog"
	"os"

	"github.com/go-chi/httplog/v2"
)

type mode int

const (
	DEV = iota
	PROD
)

var LogLevel = new(slog.LevelVar)

type Options struct {
	Logger *slog.Logger
	Port   uint16
	Mode   mode

	// Storage Options
	Bucket           string
	StorageEndpoint  string
	StorageAccessKey string
	StorageSecretKey string

	// Database Options
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     uint16
	DBDatabase string
}

func Default() *Options {
	logger := slog.New(httplog.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{
		Level: LogLevel,
	}))

	return &Options{
		Logger: logger,
		Port:   8080,
		Mode:   PROD,

		// Postgres defaults
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBHost:     "127.0.0.1",
		DBPort:     5432,
		DBDatabase: "postgres",
	}
}
