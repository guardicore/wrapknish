package http

// wrapped package
import "net/http"

import (
	"io"
	"net/url"
)

// WARNING: unsupported variable type *ast.ValueSpec for &{Doc:<nil> Names:[DefaultClient] Type:<nil> Values:[0xc42011a160] Comment:<nil>}
var DefaultClient = http.DefaultClient

type RoundTripper = http.RoundTripper

func Get(url string) (resp *Response, err error) { return http.Get(url) }
func (_recv_o *Client) Get(url string) (resp *Response, err error) {
	return (*http.Client).Get((*http.Client)(_recv_o), url)
}

var ErrUseLastResponse = http.ErrUseLastResponse

func (_recv_o *Client) Do(req *Request) (*Response, error) {
	return (*http.Client).Do((*http.Client)(_recv_o), req)
}
func (_recv_o *Client) Post(url string, contentType string, body io.Reader) (resp *Response, err error) {
	return (*http.Client).Post((*http.Client)(_recv_o), url, contentType, body)
}
func PostForm(url string, data url.Values) (resp *Response, err error) {
	return http.PostForm(url, data)
}
func (_recv_o *Client) PostForm(url string, data url.Values) (resp *Response, err error) {
	return (*http.Client).PostForm((*http.Client)(_recv_o), url, data)
}
func Head(url string) (resp *Response, err error) { return http.Head(url) }
