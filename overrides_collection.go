package main

import (
    "path/filepath"
    "log"
    "go/ast"
    "go/token"
    "go/parser"
)

func collectGenDeclOverrides(w *WrapperOverrides, d *ast.GenDecl) {
    switch d.Tok {
    case token.IMPORT:
        ;
    case token.TYPE:
        for _, s := range d.Specs {
            typeSpec := s.(*ast.TypeSpec)
            if ast.IsExported(typeSpec.Name.Name) {
                w.MarkTypeOverride(typeSpec.Name.Name, TYPE_OVERRIDE_AUTO_WRAP_METHODS)
                w.MarkTypeOverride("*" + typeSpec.Name.Name, TYPE_OVERRIDE_AUTO_WRAP_METHODS)
                log.Printf("marked type name %s for override\n", typeSpec.Name.Name)
            }
        }
    case token.VAR:
        fallthrough
    case token.CONST:
        for _, s := range d.Specs {
            valSpec := s.(*ast.ValueSpec)
            for _, n := range valSpec.Names {
                if ast.IsExported(n.Name) {
                    w.MarkVarOverride(n.Name)
                    log.Printf("marked var name %s for override\n", n.Name)
                }
            }
        }
    default:
        log.Printf("encountered unexpected type %d (%T): %+v\n", d.Tok, d, d)
        fatal("gendecl has invalid type")
    }
}

func collectFuncOverrides(w *WrapperOverrides, f *ast.FuncDecl) {
    if !ast.IsExported(f.Name.Name) {
        return
    }
    recvTypeName := getFuncRecvType(f)
    if recvTypeName == "" {
        log.Printf("marked func name %s for override\n", f.Name.Name)
        w.MarkFuncOverride(f.Name.Name)
    } else {
        if recvTypeName[0] == '*' {
            recvTypeName = recvTypeName[1:]
        }
        if ast.IsExported(recvTypeName) {
            log.Printf("marked method name %s.%s() for override\n", recvTypeName, f.Name.Name)
            w.MarkMethodOverride(recvTypeName, f.Name.Name)
        }
    }
}

func collectOverridesFromFile(w *WrapperOverrides, f *ast.File) {
    for _, decl := range f.Decls {
        switch v := decl.(type) {
        case *ast.GenDecl:
            collectGenDeclOverrides(w, v)
        case *ast.FuncDecl:
            collectFuncOverrides(w, v)
        default:
            log.Printf("unsupported declaration of type %T: %+v\n", v, v)
            fatal("unsupported declaration type")
        }
    }
}

func collectOverrides(w *WrapperOverrides, pkg string, dir string) (*token.FileSet) {
    fset := token.NewFileSet()
    finalDir := filepath.Join(dir, pkg)
    pmap, err := parser.ParseDir(fset, finalDir, nil, 0)
    if err != nil {
        fatal(err.Error())
    }
    pkgBase := filepath.Base(pkg)
    for pkgname, p := range pmap {
        if pkgname != pkgBase {
            log.Printf("skipping override package %s\n", pkgname)
            continue
        }

        for filename, f := range p.Files {
            log.Printf("collecting override definitions from %s\n", filename)
            collectOverridesFromFile(w, f)
        }
    }
    return fset
}

