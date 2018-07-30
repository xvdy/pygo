package vm

import (
	"fmt"
	"log"
	"pygo/vm/types"
)

type CodeFlags struct {
	optimized,
	newlocals,
	varargs,
	varkeywords,
	nested,
	generator uint32
}

type PyCode struct {
	types.PyObjectData

	argcount,
	nlocals,
	stacksize,
	raw_flags,
	firstlineno uint32

	consts,
	names,
	varnames,
	freevars,
	cellvars,
	lnotab types.PyObject

	filename,
	name *string

	flags CodeFlags

	vm   *VM
	code *codeReader
}

func (pc *PyCode) GetFromCellFreeStorage(idx int) (obj types.PyObject) {
	// TODO: Add some range checks

	if idx < len(pc.cellvars.(*types.PyTuple).Items) {
		// variable is cellvars[i] if i is less than the length of cellvars
		obj = pc.cellvars.(*types.PyTuple).Items[idx]
	} else {
		obj = pc.freevars.(*types.PyTuple).Items[idx-len(pc.cellvars.(*types.PyTuple).Items)]
	}
	return
}

func (pc *PyCode) SetCellFreeStorage(idx int, item types.PyObject) error {

	return nil
}

func (pc *PyCode) Log(msg string, debug bool) {
	if debug && !pc.vm.debug {
		// Ignore debug messages if debug mode is off
		return
	}
	if debug {
		log.Println(fmt.Sprintf("[%s/%s:DEBUG] %s", *pc.filename, *pc.name, msg))
	} else {
		log.Println(fmt.Sprintf("[%s/%s] %s", *pc.filename, *pc.name, msg))
	}
}

func (pc *PyCode) runtimeError(msg string) {
	panic(fmt.Sprintf("[%s/%s] Runtime error: %s", *pc.filename, *pc.name, msg))
}

func (pc *PyCode) AsString() *string {
	return pc.name
}

func (pc *PyCode) GetValue() interface{} {
	return pc.code
}

func (pc *PyCode) IsTrue() bool {
	return true
}

func (pc *PyCode) GetAttr(name, standard types.PyObject) types.PyObject {
	return nil
}

func NewPyCode(vm *VM) types.PyObject {
	co := new(PyCode)
	co.PyObjInit()
	co.vm = vm
	co.argcount, _ = vm.content.readDWord()
	co.nlocals, _ = vm.content.readDWord()
	co.stacksize, _ = vm.content.readDWord()
	co.raw_flags, _ = vm.content.readDWord()
	co.code = NewCodeReader([]byte(*vm.readObject().GetValue().(*string)))
	co.consts = vm.readObject()
	co.names = vm.readObject()
	co.varnames = vm.readObject()
	co.freevars = vm.readObject()
	co.cellvars = vm.readObject()
	co.filename = vm.readObject().GetValue().(*string)
	co.name = vm.readObject().GetValue().(*string)
	co.firstlineno, _ = vm.content.readDWord()
	co.lnotab = vm.readObject()
	co.parseFlags()
	return co
}

func (co *PyCode) parseFlags() {
	co.flags.optimized = co.raw_flags & 0x0001
	co.flags.newlocals = co.raw_flags & 0x0002
	co.flags.varargs = co.raw_flags & 0x0004
	co.flags.varkeywords = co.raw_flags & 0x0008
	co.flags.nested = co.raw_flags & 0x0010
	co.flags.generator = co.raw_flags & 0x0020
}
