package main

import (
    "log"
    "strings"
    "bytes"
    "go/ast"
    "go/token"
)

func isErrorVariableName(name string) (bool) {
    return strings.HasPrefix(name, "Err") || strings.HasSuffix(name, "Err")
}

func isErrorPointer(vspec *ast.ValueSpec) (bool) {
    for i, val := range vspec.Values {
        if i >= len(vspec.Names) {
            continue
        }
        name := vspec.Names[i].Name
        if !ast.IsExported(name) {
            continue
        }
        if !isErrorVariableName(name) {
            log.Printf("variable name %s is not an error name\n", name)
            return false
        }

        switch v := val.(type) {
        case *ast.UnaryExpr:
            if v.Op != token.AND {
                log.Printf("UnaryExpr %+v is not an addressing operation\n", v)
                return false
            }

            composite, ok := v.X.(*ast.CompositeLit)
            if !ok {
                log.Printf("UnaryExpr %+v does not address an CompositeLit object\n", v)
                return false
            }

            typeNameBuffer := new(bytes.Buffer)
            emitType(typeNameBuffer, composite.Type)
            typeName := typeNameBuffer.String()
            if !strings.HasSuffix(typeName, "Error") {
                log.Printf("CompositeLit %+v does not use an error type (%s)\n", v, typeName)
                return false
            }

        case *ast.SelectorExpr:
            if !isErrorVariableName(v.Sel.Name) {
                log.Printf("SelectorExpr %+v does not select an error variable (%s)\n", v, v.Sel.Name)
                return false
            }

        case *ast.CallExpr:
            callSelector, ok := v.Fun.(*ast.SelectorExpr)
            if !ok {
                log.Printf("CallExpr %+v does not have a selector expression\n", v)
                return false
            }
            pkg, ok := callSelector.X.(*ast.Ident)
            if !ok {
                log.Printf("SelectorExpr %+v does not have an Ident expression\n", callSelector)
                return false
            }

            if pkg.Name != "errors" && pkg.Name != "awserr" {
                log.Printf("package %s in variable value assignment is not \"errors\" or \"awserr\"\n", pkg.Name)
                return false
            }

            if callSelector.Sel.Name != "New" {
                log.Printf("function %s in variable value assignment is not \"errors.New()\"\n", callSelector.Sel.Name)
                return false
            }

        default:
            log.Printf("ValueSpec %+v of type %T is not an error pointer\n", v, v)
            return false
        }
    }
    return true
}

func isWhitelistedGlobal(vspec *ast.ValueSpec) (bool) {
    return false
}

func isAllowedGlobalPointer(vspec *ast.ValueSpec) (bool) {
    if isErrorPointer(vspec) {
        return true
    }
    if isWhitelistedGlobal(vspec) {
        return true
    }

    return false
}

