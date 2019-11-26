package http
// ^ Make sure to give your override code the same package name as the code you're overriding.

// Import anything you need.
import (
    "io"
    "fmt"
    "net/http"
)

// The regular function case is simple -- just give it your own definition that matches the original signature.
func Post(url string, contentType string, body io.Reader) (resp *Response, err error) {
    fmt.Printf("inside wrapped function! about to POST to %s\n", url)
    // You can call the original function by specifying its package (remember that within the package code you don't do this,
    // so we won't have a conflict here):
    r, e := http.Post(url, contentType, body)
    fmt.Printf("inside wrapped function! POST completed with error value %v\n", e)
    return r, e
}

// Overriding a method is a bit more complicated. First, you need to create a new type with an identical layout to the original:
type Client http.Client
// Note that you can't use alias syntax (type Client = http.Client) here, since that wouldn't enable us to redefine methods.

// Now for the method itself:
func (c *Client) Head(url string) (resp *Response, err error) {
    fmt.Printf("inside wrapped method! about to client.HEAD to %s\n", url)

    // Things start to get ugly when you want to call the original version of the method.
    // First, get the "function" version of that method.
    f := (*http.Client).Head

    // Now force our receiver pointer in there:
    r, e := f((*http.Client)(c), url)

    fmt.Printf("inside wrapped method! client.HEAD completed with error value %v\n", e)

    return r, e
}

// The rest of http.Client's methods will automatically be wrapped, so you don't have to worry about implementing everything.

