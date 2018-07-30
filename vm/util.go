package vm

import (
	"pygo/vm/types"
)

func NewPyFuncInternal(module *types.Module, fn types.ModuleFunc, name *string) types.PyObject {
	return types.NewPyFunc(types.PyFuncInternal, module, fn, name, nil)
}

func NewPyFuncExternal(codeobj types.PyObject) types.PyObject {
	return types.NewPyFunc(types.PyFuncExternal, nil, nil,nil, codeobj)
}
