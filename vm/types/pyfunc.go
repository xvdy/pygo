package types

import (
	"fmt"
	"log"
)

const (
	PyFuncExternal = iota
	PyFuncInternal
)

type PyFunc struct {
	PyObjectData

	// Internal func
	module *Module
	name   *string
	mfunc  ModuleFunc

	// External func
	codeobj PyObject

	// Meta
	functype int
	closure  PyObject // nil or tuple of cell objects
	_doc     PyObject // __doc__ attribute | not used yet
	_name    PyObject // __name__ attribute | not used yet
}

func (pf *PyFunc) SetClosure(c PyObject) {
	pf.closure = c
}

func (pf *PyFunc) GetValue() interface{} {
	return pf.codeobj
}

func (pf *PyFunc) IsTrue() bool {
	return true
}

func (pf *PyFunc) IsExternal() bool {
	return pf.functype == PyFuncExternal
}

func (pf *PyFunc) AsString() *string {
	var str string
	switch pf.functype {
	case PyFuncInternal:
		str = fmt.Sprintf("<internal function %s>", *pf.name)
	case PyFuncExternal:
		str = fmt.Sprintf("<external function %s>", *pf.codeobj.AsString())
	default:
		panic("unknown func type")
	}
	return &str
}

func (pf *PyFunc) log(msg string) {
	var ident string

	switch pf.functype {
	case PyFuncInternal:
		ident = fmt.Sprintf("%s.%s", pf.module.name, *pf.name)
	case PyFuncExternal:
		panic("not implement")
		//ident = fmt.Sprintf("%s/%s", *pf.codeobj.(*vm.PyCode).filename, *pf.codeobj.(*types.PyCode).Name)
	default:
		panic("unknown func type")
	}

	log.Println(fmt.Sprintf("[%s] %s", ident, msg))
}

func (pf *PyFunc) Run(args *PyArgs) PyObject {
	// Create frame for run
	//frame := vm.NewPyFrame(1000) // TODO change stack size to a better value?

	//starttime := time.Now()
	//defer func() {
	//	if vm.DebugMode {
	//		pf.log(fmt.Sprintf("Execution took %s.", time.Since(starttime)))
	//	}
	//}()

	switch pf.functype {
	case PyFuncInternal:
		return pf.mfunc(args)
	case PyFuncExternal:
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

func NewPyFunc(funcType int, module *Module, fn ModuleFunc, name *string, codeobj PyObject) PyObject {
	if funcType == PyFuncInternal {
		return &PyFunc{
			module:   module,
			name:     name,
			functype: PyFuncInternal,
			mfunc:    fn,
		}
	} else {
		pf := &PyFunc{
			codeobj:  codeobj,
			functype: PyFuncExternal,
		}
		pf.PyObjInit()
		return pf
	}
}

func NewFuncInternal(module *Module, fn ModuleFunc, name *string) PyObject {
	return &PyFunc{
		module:   module,
		name:     name,
		functype: PyFuncInternal,
		mfunc:    fn,
	}
}

func NewFuncExternal(codeobj PyObject) PyObject {
	pf := &PyFunc{
		codeobj:  codeobj,
		functype: PyFuncExternal,
	}
	pf.PyObjInit()

	return pf
}
