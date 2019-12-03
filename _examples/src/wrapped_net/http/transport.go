package http

// wrapped package
import "net/http"

import (
	"net/url"
)

// WARNING: unsupported variable type *ast.ValueSpec for &{Doc:<nil> Names:[DefaultTransport] Type:RoundTripper Values:[0xc42067a180] Comment:<nil>}
var DefaultTransport = http.DefaultTransport

const DefaultMaxIdleConnsPerHost = 2

type Transport = http.Transport

func ProxyFromEnvironment(req *Request) (*url.URL, error)         { return http.ProxyFromEnvironment(req) }
func ProxyURL(fixedURL *url.URL) func(*Request) (*url.URL, error) { return http.ProxyURL(fixedURL) }

var ErrSkipAltProtocol = http.ErrSkipAltProtocol
