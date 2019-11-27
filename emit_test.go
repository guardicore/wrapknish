package main

import (
    "testing"
    "bytes"
    "log"
    "go/parser"
    "go/ast"
    "go/token"
)

func getAstDecls(decls string) []ast.Decl {
    fset := token.NewFileSet() // positions are relative to fset
    src := "package mypkg\n\n" + decls
    f, err := parser.ParseFile(fset, "src.go", src, 0)
    if err != nil {
        log.Printf("failed to parse test file: %s", err.Error())
        panic("parse failed")
    }

    return f.Decls
}

func getAstDecl(decl string) ast.Decl {
    decls := getAstDecls(decl)
    if len(decls) != 1 {
        panic("incorrect number of decls: " + string(len(decls)))
    }
    return decls[0]
}

func TestEmitConstSpecs(t *testing.T) {
    valueSpec :=
`const (
    ExportedConst = 1234
    AnotherExportedConst = 6789
    unexportedConst = 9999
    OverriddenConst = 0
)`
    expectedSpec :=
`const (
    ExportedConst = 1234
    AnotherExportedConst = 6789
)
`

    overrides := NewWrapperOverrides()
    overrides.MarkVarOverride("OverriddenConst")
    constDecl := getAstDecl(valueSpec).(*ast.GenDecl)
    constBuffer := new(bytes.Buffer)
    emitConstSpecs(constBuffer, constDecl.Specs, overrides)
    emittedConst := constBuffer.String()
    if emittedConst != expectedSpec {
        t.Errorf("emitted const spec was incorrect: expected\n\n\"\"\"%s\"\"\"\n\nbut got\n\n\"\"\"%s\"\"\"\n", expectedSpec, emittedConst)
    }
}

func TestEmitValueSpecs(t *testing.T) {
    valueSpec :=
`var (
    SomeExportedVar = 1234
    someUnexportedVar = 6789
    JustAnErr = errors.New("something terrible")
    SomeOverriddenVar = 0
)`
    expectedSpec :=
`// WARNING: unsupported variable type *ast.ValueSpec for [SomeExportedVar]
var SomeExportedVar = orig_pkg.SomeExportedVar
var JustAnErr = orig_pkg.JustAnErr
`

    overrides := NewWrapperOverrides()
    overrides.MarkVarOverride("SomeOverriddenVar")
    varDecl := getAstDecl(valueSpec).(*ast.GenDecl)
    varBuffer := new(bytes.Buffer)
    emitValueSpecs(varBuffer, varDecl.Specs, overrides, "orig_pkg")
    emittedVar := varBuffer.String()
    if emittedVar != expectedSpec {
        t.Errorf("emitted var spec was incorrect: expected\n\n\"\"\"%s\"\"\"\n\nbut got\n\n\"\"\"%s\"\"\"\n", expectedSpec, emittedVar)
    }
}

func TestEmitTypeSpecs(t *testing.T) {
    typeSpec :=
`type (
    MyType struct {
        thing1 int
        thing2 string
        thing3 char
    }

    OverriddenType struct {
        thing1 int
        thing2 string
        thing3 char
    }

    unexportedType struct {
        thing1 int
        thing2 string
        thing3 char
    }

    AnInterface interface {}
)`
    expectedSpec :=
`type MyType = orig_pkg.MyType
type AnInterface = orig_pkg.AnInterface

`

    overrides := NewWrapperOverrides()
    overrides.MarkTypeOverride("OverriddenType", TYPE_OVERRIDE_AUTO_WRAP_METHODS)
    typeDecl := getAstDecl(typeSpec).(*ast.GenDecl)
    typeBuffer := new(bytes.Buffer)
    emitTypeSpecs(typeBuffer, typeDecl.Specs, overrides, "orig_pkg")
    emittedType := typeBuffer.String()
    if emittedType != expectedSpec {
        t.Errorf("emitted type spec was incorrect: expected\n\n\"\"\"%s\"\"\"\n\nbut got\n\n\"\"\"%s\"\"\"\n", expectedSpec, emittedType)
    }
}

func TestEmitImportSpecs(t *testing.T) {
    importSpec :=
`import (
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

func TestEmitFuncWrapper(t *testing.T) {
    expectedValues := []struct{
        funcDecl string
        output string
    }{
        { "func RegularExportedFunc(param1 string, param2 int) (char, error) { }", "func RegularExportedFunc(param1 string, param2 int) (char, error) { return orig_pkg.RegularExportedFunc(param1, param2) }\n" },
        { "func unexportedFunc(param1 char, param2 int) string { }", "" },
        { "func VariadicFunc(param1 string, moreParams ...int) error { }", "func VariadicFunc(param1 string, moreParams ...int) error { return orig_pkg.VariadicFunc(param1, moreParams...) }\n" },
        { "func NoReturns() { }", "func NoReturns() { orig_pkg.NoReturns() }\n" },
        { "func OverriddenFunc(param1 char) (string, error, int)", "" },
        { "func AnnoyingTypeSyntax(param1, param2, param3 string, param4 char, param5, param6, param7 int) (return1, return2 int, return3 string) { }", "func AnnoyingTypeSyntax(param1, param2, param3 string, param4 char, param5, param6, param7 int) (return1, return2 int, return3 string) { return orig_pkg.AnnoyingTypeSyntax(param1, param2, param3, param4, param5, param6, param7) }\n" },
        { "func MissingParameterNames(string, char, int) { }", "func MissingParameterNames(param_0 string, param_1 char, param_2 int) { orig_pkg.MissingParameterNames(param_0, param_1, param_2) }\n" },
        { "func (c *MyType) SomeModifyingMethod(param1 string, otherParams ...int) (string, string, string) { }", "func (_recv_o *MyType) SomeModifyingMethod(param1 string, otherParams ...int) (string, string, string) { return (*orig_pkg.MyType).SomeModifyingMethod((*orig_pkg.MyType)(_recv_o), param1, otherParams...) }\n" },
        { "func (c MyType) SomeNonModifyingMethod(string, int) { }", "func (_recv_o MyType) SomeNonModifyingMethod(param_0 string, param_1 int) { (*orig_pkg.MyType).SomeNonModifyingMethod((*orig_pkg.MyType)(_recv_o), param_0, param_1) }\n" },
    }

    overrides := NewWrapperOverrides()
    overrides.MarkFuncOverride("OverriddenFunc")
    overrides.MarkTypeOverride("MyType", TYPE_OVERRIDE_AUTO_WRAP_METHODS)
    overrides.MarkMethodOverride("MyType", "OverriddenMethod")

    for _, v := range expectedValues {
        funcDecl := getAstDecl(v.funcDecl).(*ast.FuncDecl)
        funcBuffer := new(bytes.Buffer)
        emitFuncWrapper(funcBuffer, funcDecl, overrides, "orig_pkg")
        emittedFunc := funcBuffer.String()

        if emittedFunc != v.output {
            t.Errorf("emitted wrapper function was incorrect: expected\n\n\"\"\"%s\"\"\"\n\nbut got\n\n\"\"\"%s\"\"\"\n", v.output, emittedFunc)
        }
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

