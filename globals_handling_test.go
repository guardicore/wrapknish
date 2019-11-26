package main

import (
    "testing"
    "go/ast"
)

func TestIsErrorVariableName(t *testing.T) {
    verdicts := []struct{
        name string
        verdict bool
    }{
        { "ErrSomethingWentWrong", true },
        { "SomeSuccessfulThingHappened", false },
        { "TryHarderNextTimerr", false },
        { "errolSpence", false },
        { "ShoutingIntoStderr", false },
        { "OutOfMemoryErr", true},
    }

    for _, v := range verdicts {
        if isErrorVariableName(v.name) != v.verdict {
            t.Errorf("incorrect verdict for %q", v.name)
        }
    }
}

func TestIsErrorPointer(t *testing.T) {
    verdicts := []struct{
        expr string
        verdict bool
    }{
        { "var NotAnErrThisTime = 1234", false },
        { "var ErrUnexpectedEOF = errors.New(\"unexpected EOF\")", true },
        { "var MyErr = &SyntaxError{\"invalid boolean\"}", true },
        { "var NotReallyAnError = errors.Old(\"everything is okay\")", false },
        { "var SomeoneElseErr = io.ErrUnexpectedEOF", true },
        { "var PretendingToBeAnErr = 1234", false },
    }
    for _, v := range verdicts {
        decl := getAstDecl(v.expr).(*ast.GenDecl)
        spec := decl.Specs[0].(*ast.ValueSpec)
        if isErrorPointer(spec) != v.verdict {
            t.Errorf("incorrect verdict for %q", v.expr)
        }
    }
}
