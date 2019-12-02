# wrapknish go package wrapper generator

wrapknish automatically generates wrappers for Go packages that expose an identical API but enable you to add your own custom functionality wherever you need it. This is helpful for transparent code instrumentation, stubbing for tests, and any number of other uses. For each exported type, function, variable, or constant definition in the chosen package, wrapknish creates a new one with the same name. 

For (non-method) functions, the generated function simply calls the original one with the same parameters. 

For constants (and enums), an identical value is exported. 

Global variables are more problematic, since typically you will want to preserve the original library's semantics, and Go doesn't have variable-aliasing functionality at the moment (see the proposal [here](https://go.googlesource.com/proposal/+/1487446b91599daa695905dc51a77d1bcc7086d8/design/16339-alias-decls.md)). Fortunately, in many cases, globals in real-life libraries are essentially constant pointers (as Go doesn't actually have support for const pointers). In particular, you see this for error objects, which rarely change after creation. wrapknish heuristically identifies error globals and exports identical pointers. For all other globals, it emits a warning.

Any type that hasn't explicitly been overridden is aliased (i.e., `type MyType = somepkg.MyType`). Types that have been overridden (typically with `type MyType somepkg.MyType` to preserve the original struct layout) will have wrapper methods generated for each method that hasn't explicitly been overridden. This is done with some ugly casting that relies on the type keeping its layout (which isn't strictly true but is the case in reality).

You can override any function (or method), variable, constant, type, or even an entire source file by pointing wrapknish to an overrides directory that contains any declarations you want to omit from the original package. Any name that appears in the overrides directory will be dropped in the generated output, and all files in the overrides directory will be copied into the resulting package directory. (This allows you to override files by simply creating an override file of the same name.)

## Requirements
In order to clean up unnecessary imports from generated Go files, wrapknish uses the [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) command.
```sh
apt install golang-golang-x-tools
```
Other than goimports, wrapknish uses only standard go packages.

## Building
```sh
go build
```

## Usage
```sh
wrapknish -p net/http -n hooked_ -c /home/user/my_http_hooks/
```
will generate a package under `hooked_net/http` in the current directory. If you want to override the http.Get() function, you can create a file named `/home/user/my_http_hooks/net/http/my_get.go` with your definition:
```go
package http

import (
  "fmt"
  "http"
)

func Get(url string) (resp *Response, err error) {
  fmt.Print("performing get on %s\n", url)
  return http.Get(url)
}
```

Command line flags:
```
-p <package>
    the package to wrap (for example, "net/http").
    wrapknish will search for this package under $GOROOT and $GOPATH.
-n <prefix>
    prefix for new package path.
    for example, "hooked_" will create "hooked_net/http" from "net/http".
-c <overrides>
    directory of code to override wrappers with.
    for example, "my_hooks" will have override implementations under "my_hooks/net/http/my_impl.go" if you wrap "net/http".
```
