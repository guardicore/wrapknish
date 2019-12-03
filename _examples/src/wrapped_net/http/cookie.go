package http

// wrapped package
import "net/http"

type Cookie = http.Cookie

func SetCookie(w ResponseWriter, cookie *Cookie) { http.SetCookie(w, cookie) }
