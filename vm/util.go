package vm

import (
	"pygo/vm/types"
)

func NewPyFuncInternal(module *types.Module, fn types.ModuleFunc, name *string) types.PyObject {
	return types.NewPyFunc(types.PyFuncInternal, module, fn, name, nil)
}

func NewPyFuncExternal(codeobj types.PyObject) types.PyObject {
	return types.NewPyFunc(types.PyFuncExternal, nil, nil, nil, codeobj)
}

func RunPyFunc(pf *types.PyFunc, argpes *types.PyArgs) types.PyObject {
	// Create frame for run
	//frame := NewPyFrame(1000) // TODO change stack size to a better value?
	//
	//starttime := time.Now()
	//defer func() {
	//	if DebugMode{
	//		pf.log(fmt.sprintf("execution took %s.", time.since(starttime)))
	//	}
	//}(starttime)

	switch pf.GetFuncType() {
	case types.PyFuncInternal:
		return pf.mfunc(args)
	case types.PyFuncExternal:
		panic("not implement")
		//if args != nil {
		//	for i, value := range args.positional {
		//		// Iterate reverse! Therefore:
		//		idx := len(args.positional) - 1 - i
		//
		//		name := pf.codeobj.(*PyCode).varnames.(*PyTuple).Items[idx]
		//		frame.Names[*name.AsString()] = value
		//		//fmt.Printf("\n  --- Setting %v -> %v...\n", *name.asString(), *value.asString())
		//	}
		//	if len(args.keyword) > 0 {
		//		panic("Not implemented")
		//	}
		//}
		//if vm.DebugMode {
		//	pf.log("Called")
		//}
		//res, err := pf.codeobj.(*vm.PyCode).Eval(frame)
		//if err != nil {
		//	pf.codeobj.(*vm.PyCode).runtimeError(err.Error())
		//}
		//return res
	default:
		panic("unknown func type")
	}
	panic("unreachable")
}
