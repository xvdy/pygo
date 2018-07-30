package types

import (
	"fmt"
	"log"
)

type ModuleFunc func(args *PyArgs) PyObject
type ModuleDict map[string]ModuleFunc

type Module struct {
	name  string
	funcs ModuleDict
}

func (m *Module) inject(pm *PyModule) error {
	for name, fn := range m.funcs {
		log.Printf("Injecting %v -> %v\n", name, fn)
		//pm.attributes[name] = NewPyFuncInternal(m, fn, &name)
		pm.attributes[name] = NewPyFunc(PyFuncInternal, m, fn, &name, nil)
	}
	return nil
}

type PyModule struct {
	PyObjectData
	module *Module
	name   *string

	// If module is extern (it's own pyc file):
	//content          *vm.codeReader
	code             PyObject
	interned_strings []PyObject
}

func (pm *PyModule) AsString() *string {
	str := fmt.Sprintf("<module %s>", *pm.name)
	return &str
}

func (pm *PyModule) GetValue() interface{} {
	return nil
}

func (pm *PyModule) IsTrue() bool {
	return true
}

func NewPyModule(name *string) PyObject {
	mod := new(PyModule)
	mod.PyObjInit()
	mod.name = name

	// Import all functions and global names and make them
	// available in the attributes
	module, is_builtin := Modules[*name]
	if is_builtin {
		mod.module = &module
		if err := module.inject(mod); err != nil {
			panic("Error during module injection: " + err.Error())
		}
	} else {
		// Search for a pyc file and execute it!
		panic(fmt.Sprintf("Non-builtin modules are not supported yet (%v)", *name))
	}

	return mod
}
