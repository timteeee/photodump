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
	"photodump/internal/options"
	"text/template"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type App struct {
	logger      *slog.Logger
	port        uint16
	objectStore *minio.Client
	bucket      string
	dbUser      string
	dbPassword  string
	dbHost      string
	dbPort      uint16
	dbDatabase  string
	migrations  fs.FS
	templates   *template.Template
	static      fs.FS
}

func (app *App) DatabaseURL() string {
	return fmt.Sprintf("pgx5://%s:%s@%s:%d/%s", app.dbUser, app.dbPassword, app.dbHost, app.dbPort, app.dbDatabase)
}

func New(opts *options.Options, templates, migrations, staticAssets fs.FS) *App {
	ctx := context.Background()

	client, err := minio.New(opts.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.StorageAccessKey, opts.StorageSecretKey, ""),
		Secure: opts.Mode == options.PROD,
	})

	if err != nil {
		fmt.Printf("Error occurred while attempting to connect to object storage at %s:\n\t%s\n", opts.StorageEndpoint, err)
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

	tmpl := template.Must(template.New("templates").Funcs(html.Funcs).ParseFS(templates, "templates/*"))

	return &App{
		port:        opts.Port,
		logger:      opts.Logger,
		templates:   tmpl,
		migrations:  migrations,
		static:      staticAssets,
		objectStore: client,
		bucket:      opts.Bucket,
		dbUser:      opts.DBUser,
		dbPassword:  opts.DBPassword,
		dbHost:      opts.DBHost,
		dbPort:      opts.DBPort,
		dbDatabase:  opts.DBDatabase,
	}
}

func (app *App) Run(ctx context.Context) error {
	if err := app.runUpMigrations(); err != nil {
		return err
	}

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

func (app *App) makeMigrator() (*migrate.Migrate, error) {
	src, err := iofs.New(app.migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("unable to create source from migrations dir:\n\t%w", err)
	}

	m, err := migrate.NewWithSourceInstance("migrations", src, app.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("unable to create migrate instance:\n\t%w", err)
	}

	return m, nil
}

func (app *App) runUpMigrations() error {
	m, err := app.makeMigrator()
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations unsuccessful:\n\t%w", err)
	}

	return nil
}

func (app *App) runDownMigrations() error {
	m, err := app.makeMigrator()
	if err != nil {
		return err
	}

	if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations unsuccessful:\n\t%w", err)
	}

	return nil
}
