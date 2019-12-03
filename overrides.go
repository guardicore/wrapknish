package main

type methodName struct {
    RecvType, FuncName string
}

type WrapperOverrides struct {
    funcs map[string]bool
    methods map[methodName]bool
    types map[string]int
    vars map[string]bool
}

const (
    TYPE_OVERRIDE_NONE int = 0
    TYPE_OVERRIDE_NO_METHODS int = 1
    TYPE_OVERRIDE_AUTO_WRAP_METHODS int = 2
)

func NewWrapperOverrides() *WrapperOverrides {
    return &WrapperOverrides {
        funcs: make(map[string]bool),
        methods: make(map[methodName]bool),
        types: make(map[string]int),
        vars: make(map[string]bool),
    }
}

func (w *WrapperOverrides) IsFuncOverridden(funcName string) bool {
    _, exists := w.funcs[funcName]
    return exists
}
func (w *WrapperOverrides) IsMethodOverridden(recvType, funcName string) bool {
    _, exists := w.methods[methodName{recvType, funcName}]
    return exists
}
func (w *WrapperOverrides) GetTypeOverride(typeName string) int {
    res, exists := w.types[typeName]
    if exists {
        return res
    } else {
        return TYPE_OVERRIDE_NONE
    }
}
func (w *WrapperOverrides) IsVarOverridden(varName string) bool {
    _, exists := w.vars[varName]
    return exists
}

func (w *WrapperOverrides) MarkFuncOverride(funcName string) {
    if w.IsFuncOverridden(funcName) {
        fatal("func %s has already been marked for override -- check for redefinition", funcName)
    }
    w.funcs[funcName] = true
}
func (w *WrapperOverrides) MarkMethodOverride(typeName, funcName string) {
    if w.IsMethodOverridden(typeName, funcName) {
        fatal("method %s has already been marked for override -- check for redefinition", funcName)
    }
    w.methods[methodName{typeName, funcName}] = true
}
func (w *WrapperOverrides) MarkTypeOverride(typeName string, overrideMode int) {
    if w.GetTypeOverride(typeName) != TYPE_OVERRIDE_NONE {
        fatal("type %s has already been marked for override -- check for redefinition", typeName)
    }
    if overrideMode != TYPE_OVERRIDE_NO_METHODS && overrideMode != TYPE_OVERRIDE_AUTO_WRAP_METHODS {
        fatal("invalid type override mode %d\n", overrideMode)
    }
    w.types[typeName] = overrideMode
}
func (w *WrapperOverrides) MarkVarOverride(varName string) {
    if w.IsVarOverridden(varName) {
        fatal("var %s has already been marked for override -- check for redefinition", varName)
    }
    w.vars[varName] = true
}

