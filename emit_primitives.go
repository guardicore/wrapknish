package main

import (
    "io"
    "log"
    "fmt"
    "os"
    "go/ast"
)

func emit(w io.Writer, format string, args ...interface{}) {
    fmt.Fprintf(w, format, args...)
    fmt.Fprintf(os.Stderr, format, args...)
}

func emitNewline(w io.Writer) {
    emit(w, "\n")
}

func emitPackageDecl(w io.Writer, pkgname string) {
    emit(w, "package %s\n", pkgname)
}

func emitImport(w io.Writer, pkgname string) {
    emit(w, "import %q\n\n", pkgname)
}

func emitType(w io.Writer, t interface{}) {
    switch v := t.(type) {
    case *ast.Ident:
        emit(w, "%s", v.Name)
    case *ast.StarExpr:
        emit(w, "*")
        emitType(w, v.X)
    case *ast.InterfaceType:
        emit(w, "interface{")
        if v.Methods.List != nil && len(v.Methods.List) != 0 {
            fatal("interface parameters with methods are not supported")
        }
        emit(w, "}")
    case *ast.SelectorExpr:
        emitValue(w, v.X)
        emit(w, ".")
        emitType(w, v.Sel)
    case *ast.ArrayType:
        if v.Len != nil {
            fatal("unsupported array len type")
        }
        emit(w, "[]")
        emitType(w, v.Elt)
    case *ast.Ellipsis:
        emit(w, "...")
        emitType(w, v.Elt)
    case *ast.MapType:
        emit(w, "map[")
        emitType(w, v.Key)
        emit(w, "]")
        emitType(w, v.Value)
    case *ast.FuncType:
        emitFuncType(w, v)
    case *ast.ChanType:
        if v.Dir == ast.SEND {
            emit(w, "<-chan ")
        } else {
            emit(w, "chan<- ")
        }
        emitType(w, v.Value)
    default:
        log.Printf("encountered unexpected type %T: %+v\n", v, v)
        fatal("expression has invalid type")
    }
}

