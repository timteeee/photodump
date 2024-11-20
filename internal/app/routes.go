package app

import (
	"fmt"
	"net/http"
	"path/filepath"
	"photodump/internal/html"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (app *App) router() http.Handler {
	r := chi.NewRouter()

	logger := httplog.Logger{
		Logger: app.logger,
		Options: httplog.Options{
			JSON: false,
		},
	}
	// register all middlewares
	r.Use(httplog.RequestLogger(&logger))

	// serve assets
	r.Mount("/static", http.FileServerFS(app.static))

	// register routes
	r.Get("/", app.Feed)
	r.Get("/photos", app.GetPhotos)
	r.Post("/photos", app.UploadPhotos)

	return r
}

func (app *App) Feed(w http.ResponseWriter, r *http.Request) {
	app.templates.ExecuteTemplate(w, "index.html", html.PhotoFeedParams{})
}

// TODO: implement
func (app *App) GetPhotos(w http.ResponseWriter, r *http.Request) {
	photos := []html.PhotoCardParams{}

	for _, p := range photos {
		app.templates.ExecuteTemplate(w, "photo-card.html", p)
	}
}

func (app *App) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, fmt.Sprintf("uh oh, %s", err), http.StatusInternalServerError)
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

	// TODO: bulk upload,
	for _, h := range headers {
		if photo, err := h.Open(); err == nil {
			defer photo.Close()

			ext := filepath.Ext(h.Filename)
			id := uuid.NewString()
			uploadName := id + ext

			app.logger.Info("uploading photo", "filename", h.Filename)
			info, err := app.objectStore.PutObject(r.Context(), app.bucket, uploadName, photo, h.Size, minio.PutObjectOptions{
				ContentType: "application/octet-stream",
			})
			if err != nil {
				app.logger.Error("error uploading photo", "error", err, "upload_name", uploadName, "filename", h.Filename, "size", h.Size)
				http.Error(w, fmt.Sprintf("failed to upload: %s", err), http.StatusInternalServerError)
				return
			}

			app.logger.Info("uploaded photo", "filename", h.Filename)

			// store original filename, bucket, key, timestamp,

			// TODO: there is probably something smarter to do than set the expiry time to 7 days
			presignedUrl, err := app.objectStore.PresignedGetObject(r.Context(), app.bucket, info.Key, time.Hour*24*7, nil)
			app.templates.ExecuteTemplate(w, "photo-card.html", html.PhotoCardParams{
				Src:        presignedUrl.String(),
				Filename:   h.Filename,
				UploadedAt: now,
			})
		}
	}
}
