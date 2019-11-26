package http

// wrapped package
import "net/http"

import (
	"bufio"
)

type Response = http.Response

var ErrNoLocation = http.ErrNoLocation

func ReadResponse(r *bufio.Reader, req *Request) (*Response, error) { return http.ReadResponse(r, req) }
