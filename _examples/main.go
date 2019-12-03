package main

import (
    "fmt"
    "time"
    original_http "net/http"
    "wrapped_net/http"
)

func main() {
    timeout := time.Duration(10 * time.Second)
    fmt.Printf("client.Head() with wrapped package:\n")
    client := http.Client {
        Timeout: timeout,
        Transport: http.DefaultTransport,
    }
    resp, err := client.Head("https://www.google.com")
    if err != nil {
        panic("HEAD failed with wrapped package")
    }
    resp.Body.Close()
    fmt.Printf("client.Head() with wrapped package complete.\n\n")


    fmt.Printf("client.Head() with original package:\n")
    wrapped_client := original_http.Client {
        Timeout: timeout,
        Transport: http.DefaultTransport,
    }
    resp, err = wrapped_client.Head("https://www.google.com")
    if err != nil {
        panic("HEAD failed with original package")
    }
    resp.Body.Close()
    fmt.Printf("client.Head() with original package complete.\n\n")
}

