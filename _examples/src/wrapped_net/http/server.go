package http

// wrapped package
import "net/http"

import (
	"net"
	"time"
)

var ErrBodyNotAllowed = http.ErrBodyNotAllowed
var ErrHijacked = http.ErrHijacked
var ErrContentLength = http.ErrContentLength
var ErrWriteAfterFlush = http.ErrWriteAfterFlush

type Handler = http.Handler

type ResponseWriter = http.ResponseWriter

type Flusher = http.Flusher

type Hijacker = http.Hijacker

type CloseNotifier = http.CloseNotifier

// WARNING: unsupported variable type *ast.ValueSpec for &{Doc:<nil> Names:[ServerContextKey] Type:<nil> Values:[0xc42050f780] Comment:<nil>}
var ServerContextKey = http.ServerContextKey

// WARNING: unsupported variable type *ast.ValueSpec for &{Doc:<nil> Names:[LocalAddrContextKey] Type:<nil> Values:[0xc42050f800] Comment:<nil>}
var LocalAddrContextKey = http.LocalAddrContextKey

const TrailerPrefix = "Trailer:"
const DefaultMaxHeaderBytes = 1 << 20
const TimeFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

var ErrAbortHandler = http.ErrAbortHandler

type HandlerFunc = http.HandlerFunc

func Error(w ResponseWriter, error string, code int)              { http.Error(w, error, code) }
func NotFound(w ResponseWriter, r *Request)                       { http.NotFound(w, r) }
func NotFoundHandler() Handler                                    { return http.NotFoundHandler() }
func StripPrefix(prefix string, h Handler) Handler                { return http.StripPrefix(prefix, h) }
func Redirect(w ResponseWriter, r *Request, url string, code int) { http.Redirect(w, r, url, code) }
func RedirectHandler(url string, code int) Handler                { return http.RedirectHandler(url, code) }

type ServeMux = http.ServeMux

func NewServeMux() *ServeMux { return http.NewServeMux() }

// WARNING: unsupported variable type *ast.ValueSpec for &{Doc:<nil> Names:[DefaultServeMux] Type:<nil> Values:[0xc4205c0ba0] Comment:<nil>}
var DefaultServeMux = http.DefaultServeMux

func Handle(pattern string, handler Handler) { http.Handle(pattern, handler) }
func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	http.HandleFunc(pattern, handler)
}
func Serve(l net.Listener, handler Handler) error { return http.Serve(l, handler) }
func ServeTLS(l net.Listener, handler Handler, certFile, keyFile string) error {
	return http.ServeTLS(l, handler, certFile, keyFile)
}

type Server = http.Server

type ConnState = http.ConnState

const (
	StateNew ConnState = iota
	StateActive
	StateIdle
	StateHijacked
	StateClosed
)

var ErrServerClosed = http.ErrServerClosed

func ListenAndServe(addr string, handler Handler) error { return http.ListenAndServe(addr, handler) }
func ListenAndServeTLS(addr, certFile, keyFile string, handler Handler) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, handler)
}
func TimeoutHandler(h Handler, dt time.Duration, msg string) Handler {
	return http.TimeoutHandler(h, dt, msg)
}

var ErrHandlerTimeout = http.ErrHandlerTimeout
