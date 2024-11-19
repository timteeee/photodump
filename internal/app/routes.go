package app

import (
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"photodump/internal/assets"
	"photodump/internal/html"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (a *App) Router() http.Handler {
	r := chi.NewRouter()

	// register all middlewares
	r.Use(middleware.Logger)

	// serve assets
	r.Mount("/static", assets.Static())

	// register routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		page := html.HomePage{}
		page.Render(w)
	})

	// TODO: implement
	r.Get("/photos", func(w http.ResponseWriter, r *http.Request) {
		photos := []html.PhotoCard{}

		for _, p := range photos {
			p.Render(w)
		}
	})

	r.Post("/photos", func(w http.ResponseWriter, r *http.Request) {
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

		for _, h := range headers {
			if photo, err := h.Open(); err == nil {
				defer photo.Close()

				ext := filepath.Ext(h.Filename)
				id := uuid.NewString()
				uploadName := id + ext

				slog.Info("uploading photo", "filename", h.Filename)
				info, err := a.objectStore.PutObject(r.Context(), a.bucket, uploadName, photo, h.Size, minio.PutObjectOptions{
					ContentType: "application/octet-stream",
				})
				if err != nil {
					slog.Error("error uploading photo", "error", err, "upload_name", uploadName, "filename", h.Filename, "size", h.Size)
					http.Error(w, fmt.Sprintf("failed to upload: %s", err), http.StatusInternalServerError)
					return
				}

				slog.Info("uploaded photo", "filename", h.Filename)

				// store original filename, bucket, key, timestamp,

				// TODO: there is probably something smarter to do than set the expiry time to 7 days
				presignedUrl, err := a.objectStore.PresignedGetObject(r.Context(), a.bucket, info.Key, time.Hour*24*7, nil)
				card := html.PhotoCard{
					Src:        presignedUrl.String(),
					Filename:   h.Filename,
					UploadedAt: now,
				}

				card.Render(w)
			}
		}
	})

	return r
}
