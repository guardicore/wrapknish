package http

// wrapped package
import "net/http"

var ErrLineTooLong = http.ErrLineTooLong
var ErrBodyReadAfterClose = http.ErrBodyReadAfterClose
