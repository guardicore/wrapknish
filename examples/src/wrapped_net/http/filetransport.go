package http

// wrapped package
import "net/http"

func NewFileTransport(fs FileSystem) RoundTripper { return http.NewFileTransport(fs) }
