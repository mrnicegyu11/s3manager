package s3manager

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/minio/minio-go/v7"
)

// HandleBucketsView renders all buckets on an HTML page.
func HandleBucketsView(s3 S3, templates fs.FS, basePath string) http.HandlerFunc {

	type pageData struct {
		BasePath string
		Buckets  []minio.BucketInfo
	}

	return func(w http.ResponseWriter, r *http.Request) {
		buckets, err := s3.ListBuckets(r.Context())
		log.Println("basePath:")
		log.Println(basePath)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error listing buckets: %w", err))
			return
		}
		data := pageData{
			BasePath: basePath,
			Buckets:  buckets,
		}

		t, err := template.ParseFS(templates, "layout.html.tmpl", "buckets.html.tmpl")
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error parsing template files: %w", err))
			return
		}
		err = t.ExecuteTemplate(w, "layout", data)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error executing template: %w", err))
			return
		}
	}
}
