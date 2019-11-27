package main

import (
    "testing"
    "go/parser"
    "go/token"
)

var SRC =
`package testpkg

import (
    "nothing/in/particular"
    "orig_pkg"
    "log"
)

func MyOverriddenFunc(param int) (string, error) {
    log.Printf("let's do this instead!")
    return "abc", nil
}

func internalFunc() {
    log.Printf("nothing to see here")
}

const (
    OverriddenConst1 = 123
    OverriddenConst2 = "abc"
    justSomePrivateConst = 999
)

var myVar = "hello"
var TheirVar = orig_pkg.TheirVar

type (
    MyOverriddenType orig_pkg.MyOverriddenType
    myInternalType struct { abc int }
)

func (me *MyType) MyOverriddenMethod() {
    log.Printf("or maybe this")
}`

func TestCollectOverrides(t *testing.T) {
    fset := token.NewFileSet() // positions are relative to fset
    f, err := parser.ParseFile(fset, "src.go", SRC, 0)
    if err != nil {
        t.Errorf("failed to parse override sample source file")
    }

    overrides := NewWrapperOverrides()
    collectOverridesFromFile(overrides, f)

    if !overrides.IsFuncOverridden("MyOverriddenFunc") {
        t.Errorf("MyOverriddenFunc was not marked as overridden")
    }
    if overrides.IsFuncOverridden("NotMyOverriddenFunc") {
        t.Errorf("NotMyOverriddenFunc was not marked as overridden")
    }
    if overrides.IsFuncOverridden("internalFunc") {
        t.Errorf("internalFunc was marked as overridden")
    }
    if overrides.IsMethodOverridden("", "MyOverriddenFunc") {
        t.Errorf("MyOverriddenFunc was marked as an overridden method instead of function")
    }
    if !overrides.IsMethodOverridden("MyType", "MyOverriddenMethod") {
        t.Errorf("MyOverriddenMethod was not marked as an overridden")
    }
    if overrides.GetTypeOverride("MyOverriddenType") != TYPE_OVERRIDE_AUTO_WRAP_METHODS {
        t.Errorf("MyOverriddenType was not marked for override with method wrapping")
    }
    if overrides.GetTypeOverride("SomeNonOverriddenType") != TYPE_OVERRIDE_NONE {
        t.Errorf("SomeNonOverriddenType was marked for override")
    }
    if overrides.GetTypeOverride("myInternalType") != TYPE_OVERRIDE_NONE {
        t.Errorf("myInternalType was marked for override")
    }
    if !overrides.IsVarOverridden("TheirVar") {
        t.Errorf("TheirVar was not overridden")
    }
    if overrides.IsVarOverridden("myVar") {
        t.Errorf("myVar was overridden")
    }
    if !overrides.IsVarOverridden("OverriddenConst1") {
        t.Errorf("OverriddenConst1 was not overridden")
    }
    if !overrides.IsVarOverridden("OverriddenConst2") {
        t.Errorf("OverriddenConst2 was not overridden")
    }
    if overrides.IsVarOverridden("justSomePrivateConst") {
        t.Errorf("justSomePrivateConst was overridden")
    }
}
