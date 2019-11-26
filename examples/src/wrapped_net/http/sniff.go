package http

// wrapped package
import "net/http"

func DetectContentType(data []byte) string { return http.DetectContentType(data) }
