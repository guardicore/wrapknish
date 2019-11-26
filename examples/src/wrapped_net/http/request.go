package http

// wrapped package
import "net/http"

import (
	"bufio"
	"io"
)

var ErrMissingFile = http.ErrMissingFile

type ProtocolError = http.ProtocolError

var ErrNotSupported = http.ErrNotSupported
var ErrUnexpectedTrailer = http.ErrUnexpectedTrailer
var ErrMissingBoundary = http.ErrMissingBoundary
var ErrNotMultipart = http.ErrNotMultipart
var ErrHeaderTooLong = http.ErrHeaderTooLong
var ErrShortBody = http.ErrShortBody
var ErrMissingContentLength = http.ErrMissingContentLength

type Request = http.Request

var ErrNoCookie = http.ErrNoCookie

func ParseHTTPVersion(vers string) (major, minor int, ok bool) { return http.ParseHTTPVersion(vers) }
func NewRequest(method, url string, body io.Reader) (*Request, error) {
	return http.NewRequest(method, url, body)
}
func ReadRequest(b *bufio.Reader) (*Request, error) { return http.ReadRequest(b) }
func MaxBytesReader(w ResponseWriter, r io.ReadCloser, n int64) io.ReadCloser {
	return http.MaxBytesReader(w, r, n)
}
