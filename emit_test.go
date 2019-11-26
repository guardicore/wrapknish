package main

import (
    "testing"
    "bytes"
    "log"
    "go/parser"
    "go/ast"
    "go/token"
)

func getAstDecl(decl string) ast.Decl {
    fset := token.NewFileSet() // positions are relative to fset
    src := "package mypkg\n\n" + decl
    f, err := parser.ParseFile(fset, "src.go", src, 0)
    if err != nil {
        log.Printf("failed to parse test file: %s", err.Error())
        panic("parse failed")
    }

    return f.Decls[0]
}

func TestEmitImportSpecs(t *testing.T) {
    importSpec := `import (
    "net/http"
    my_log "log"
    "golang.org/x/net/context/ctxhttp"
)

`
    importDecl := getAstDecl(importSpec).(*ast.GenDecl)
    importBuffer := new(bytes.Buffer)
    emitImportSpecs(importBuffer, importDecl.Specs)
    emittedImport := importBuffer.String()
    if emittedImport != importSpec {
        t.Errorf("emitted import spec was incorrect: expected\n\n\"\"\"%s\"\"\"\n\nbut got\n\n\"\"\"%s\"\"\"\n", importSpec, emittedImport)
    }
}

func TestEmitFuncType(t *testing.T) {
    funcTypes := []struct{
        signature string
        outputType string
    }{
        { "func(firstParam string, int, thirdParam char)", "func(firstParam string, int, thirdParam char)" },
        { "func(UnnamedType, ...VariadicType)", "func(UnnamedType, ...VariadicType)" },
        { "func() (RetType)", "func() RetType" },
        { "func(func(int) string) char", "func(func(int) string) char" },
    }

    for _, v := range funcTypes {
        e, err := parser.ParseExpr(v.signature)
        if err != nil {
            t.Errorf("failed to parse function signature %q", v)
        }
        typeNameBuffer := new(bytes.Buffer)
        emitType(typeNameBuffer, e)
        emittedType := typeNameBuffer.String()
        if v.outputType != emittedType {
            t.Errorf("emitted type was incorrect: expected %q but got %q", v.outputType, emittedType)
        }
    }
}

