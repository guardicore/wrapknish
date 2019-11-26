package http

// wrapped package
import "net/http"

import (
	"io"
	"time"
)

type Dir = http.Dir

type FileSystem = http.FileSystem

type File = http.File

func ServeContent(w ResponseWriter, req *Request, name string, modtime time.Time, content io.ReadSeeker) {
	http.ServeContent(w, req, name, modtime, content)
}
func ServeFile(w ResponseWriter, r *Request, name string) { http.ServeFile(w, r, name) }
func FileServer(root FileSystem) Handler                  { return http.FileServer(root) }
