package http

// wrapped package
import "net/http"

import (
	"time"
)

type Header = http.Header

func ParseTime(text string) (t time.Time, err error) { return http.ParseTime(text) }
func CanonicalHeaderKey(s string) string             { return http.CanonicalHeaderKey(s) }
