package http

// wrapped package
import "net/http"

// WARNING: unsupported variable type *ast.ValueSpec for &{Doc:<nil> Names:[NoBody] Type:<nil> Values:[0xc420507880] Comment:<nil>}
var NoBody = http.NoBody

type PushOptions = http.PushOptions

type Pusher = http.Pusher
