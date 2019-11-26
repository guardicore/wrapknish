package main

import (
    "testing"
    "bytes"
    "go/parser"
)

func TestEmitType(t *testing.T) {
    expectedValues := []string{
        "some.selector.Expr",
        "JustAType",
        "interface{}",
        "[]ArrayElements",
        // "...VariadicType", -- ParseExpr() doesn't handle this well
        "map[KeyType]ValueType",
        "<-chan ReceiveType",
        "chan<- SendType",
        "map[<-chan mypackage.SomeType]map[interface{}]chan<- []pkgpkg.Type",
    }

    for _, v := range expectedValues {
        e, err := parser.ParseExpr(v)
        if err != nil {
            t.Errorf("failed to parse type expression %q", v)
        }
        typeNameBuffer := new(bytes.Buffer)
        emitType(typeNameBuffer, e)
        emittedType := typeNameBuffer.String()
        if v != emittedType {
            t.Errorf("emitted type was incorrect: expected %q but got %q", v, emittedType)
        }
    }
}
