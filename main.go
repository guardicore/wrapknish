package main

import (
    "io"
    "strings"
    "os"
    "path/filepath"
    "flag"
    "log"
    "go/ast"
    "go/parser"
    "go/token"
    "os/exec"
    "go/build"
)

func shouldProcessGoFile(f os.FileInfo) (bool) {
    if strings.HasSuffix(f.Name(), "_test.go") {
        return false
    }
    return true
}

func dirExists(dir string) bool {
    st, err := os.Stat(dir)
    if err != nil {
        if os.IsNotExist(err) {
            return false
        } else {
            fatal(err.Error())
        }
    }
    return st.IsDir()
}

func parseEnvDirForWrap(env string, defaultVal string, subdir string, pmap map[string]*ast.Package) {
    dir, exists := os.LookupEnv(env)
    if !exists && defaultVal != "" {
        dir = defaultVal
    }
    if dir != "" {
        dir = filepath.Join(dir, subdir)
        if dirExists(dir) {
            log.Printf("searching for package in directory %s\n", dir)
            fset := token.NewFileSet()
            res, err := parser.ParseDir(fset, dir, shouldProcessGoFile, 0)
            if err != nil {
                fatal(err.Error())
            }
            for k, v := range res {
                if _, pkgExists := pmap[k]; pkgExists {
                    log.Printf("WARNING: package %s was found in multiple locations. " +
                    " Using first found occurrence. Check $GOPATH and $GOROOT for overlap.\n", k)
                } else {
                    pmap[k] = v
                }
            }
        }
    }
}

func cleanupImports(filename string) {
    goimportsPath, err := exec.LookPath("goimports")
    if err != nil {
        log.Println("WARNING: goimports was not found on your machine. " +
            "Install goimports to clean up unnecessary imports in generated wrapper files.")
    } else {
        cmd := exec.Command(goimportsPath, "-v", "-w", filename)
        e, err := cmd.CombinedOutput()

        if err != nil {
            fatal(string(e) + err.Error())
        }
    }
}

func copyFile(srcName, dstName string) {
    src, err := os.Open(srcName)
    if err != nil {
        fatal(err.Error())
    }
    defer src.Close()
    dst, err := os.Create(dstName)
    if err != nil {
        fatal(err.Error())
    }
    defer dst.Close()
    io.Copy(dst, src)
}

func main() {

    var pkg, overridesDir, newPkgPrefix string
    flag.StringVar(&pkg, "p", "", "the package to wrap (for example, \"net/http\").\nwrapknish will search for this package under $GOROOT and $GOPATH.")
    flag.StringVar(&overridesDir, "c", "", "prefix for new package path.\nfor example, \"hooked_\" will create \"hooked_net/http\" from \"net/http\".")
    flag.StringVar(&newPkgPrefix, "n", "", "directory of code to override wrappers with.\nfor example, \"my_hooks\" should have override implementations under \"my_hooks/net/http/my_impl.go\" if you wrap \"net/http\".")
    flag.Parse()

    if pkg == "" {
        flag.Usage()
        log.Println("you must specify a package name to wrap")
        os.Exit(2)
    }
    if newPkgPrefix == "" {
        flag.Usage()
        log.Println("you must specify a prefix for the new package path")
        os.Exit(2)
    }

    pkgBase := filepath.Base(pkg)

    overrides := NewWrapperOverrides()
    var overridesFset *token.FileSet
    if overridesDir != "" {
        overridesFset = collectOverrides(overrides, pkg, overridesDir)
    }

    pmap := map[string]*ast.Package{}
    pkgDir := filepath.Join("src", pkg)
    parseEnvDirForWrap("GOROOT", build.Default.GOROOT, pkgDir, pmap)
    parseEnvDirForWrap("GOPATH", build.Default.GOPATH, pkgDir, pmap)
    if len(pmap) == 0 {
        fatal("package %s was not found. Check $GOROOT or $GOPATH are defined.\n", pkg)
    }

    // i.e., xray_net/http or xray_strings
    newPkg := newPkgPrefix + pkg
    newPkgBase := filepath.Base(newPkg)

    err := os.RemoveAll(newPkg)
    if err != nil {
        fatal(err.Error())
    }
    err = os.MkdirAll(newPkg, os.ModePerm)
    if err != nil {
        fatal(err.Error())
    }

    for pkgname, p := range pmap {
        if pkgname != pkgBase {
            log.Printf("skipping package %s\n", pkgname)
            continue
        }

        for filename, f := range p.Files {
            outputFilename := filepath.Join(newPkg, filepath.Base(filename))
            log.Printf("wrapping file %q into %q\n", filename, outputFilename)

            outputFile, err := os.Create(outputFilename)
            if err != nil {
                fatal(err.Error())
            }

            emitPackageDecl(outputFile, newPkgBase)
            emitNewline(outputFile)

            emit(outputFile, "// wrapped package\n")
            emitImport(outputFile, pkg)
            emitFileWrapper(outputFile, f, overrides, pkgBase)

            if err := outputFile.Close(); err != nil {
                fatal(err.Error())
            }

            cleanupImports(outputFilename)
        }
    }

    if overridesFset != nil {
        overridesFset.Iterate(func (f *token.File) bool {
            newFilename := filepath.Join(newPkg, filepath.Base(f.Name()))
            log.Printf("copying override file from %s to %s\n", f.Name(), newFilename)
            copyFile(f.Name(), newFilename)
            return true
        })
    }
}

