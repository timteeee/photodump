package main

import (
	"fmt"
	"os"
	"photodump/internal/app"

	flag "github.com/spf13/pflag"
)

var (
	dev     *bool   = flag.Bool("dev", false, "whether to run in dev mode")
	port    *uint16 = flag.Uint16P("port", "p", 80, "port to bind to")
	storage *string = flag.StringP("storage", "s", "", "URL for object storage")
	bucket  *string = flag.StringP("bucket", "b", "", "bucket for object storage")
)

func main() {
	flag.Parse()

	accessKey := "VPP0fkoCyBZx8YU0QTjH"
	secretKey := "iFq6k8RLJw5B0faz0cKCXeQk0w9Q8UdtaFzHuw4J"

	opts := app.Options{
		Port:                 *port,
		ObjectStoreEndpoint:  *storage,
		Bucket:               *bucket,
		Secure:               !*dev,
		ObjectStoreAccessKey: accessKey,
		ObjectStoreSecretKey: secretKey,
	}

	if err := app.New(&opts).Run(); err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}
}

// func router() *chi.Mux {
// 	r := chi.NewRouter()
//
// 	r.Use(middleware.Logger)
//
// 	r.Mount("/public", assets.Public())
//
// 	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
// 		res := html.HomePage{}
// 		res.Render(w)
// 	})
//
// 	// TODO: read from db to get filenames in storage, render photo cards with pre-signed urls
// 	r.Get("/photos", func(w http.ResponseWriter, r *http.Request) {
// 		photo := filepath.Join("templates", "photo-card.html")
//
// 		for range 3 {
// 			t, err := template.ParseFiles(photo)
// 			if err != nil {
// 				http.Error(w, fmt.Sprintf("I goofed: %s", err.Error()), 500)
// 			}
// 			t.Execute(w, nil)
// 		}
// 	})
//
// 	// TODO: store photos named as hash of contents, store filenames, sizes, datetimes, etc in db, render cards with pre-signed urls
// 	// r.Post("/photos", func(w http.ResponseWriter, r *http.Request) {
// 	// 	if err := r.ParseMultipartForm(32 << 20); err != nil {
// 	// 		http.Error(w, fmt.Sprintf("uh oh, %s", err.Error()), 500)
// 	// 		return
// 	// 	}
// 	//
// 	// 	if r.MultipartForm == nil {
// 	// 		http.Error(w, "didn't receive multipart form in request", http.StatusUnprocessableEntity)
// 	// 		return
// 	// 	}
// 	//
// 	// 	headers := r.MultipartForm.File["photos"]
// 	// 	if headers == nil {
// 	// 		http.Error(w, "didn't receive photos in request", http.StatusUnprocessableEntity)
// 	// 		return
// 	// 	}
// 	//
// 	// 	now := time.Now().UTC()
// 	//
// 	// 	for _, h := range headers {
// 	// 		buf := bytes.NewBuffer(make([]byte, 0, h.Size))
// 	// 		if photo, err := h.Open(); err == nil {
// 	// 			io.Copy(buf, photo)
// 	// 			photo.Close()
// 	//
// 	// 			photoStore.Store(h.Filename, buf.Bytes())
// 	// 			slog.Info("stored photo", "filename", h.Filename)
// 	//
// 	// 			card := html.PhotoCard{
// 	// 				// TODO: create and use presigned url here
// 	// 				Src:        fmt.Sprintf("/photos/%s", h.Filename),
// 	// 				Filename:   h.Filename,
// 	// 				UploadedAt: now,
// 	// 			}
// 	// 			card.Render(w)
// 	// 		}
// 	// 	}
// 	// })
//
// 	return r
// }
