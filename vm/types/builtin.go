package types

import "time"

var PyNil = NewPyNone()
var PyTrue = NewPyBool(true)
var PyFalse = NewPyBool(false)

var PyBuiltInTypeMap = map[string]PyObject{
	"True":  PyTrue,
	"False": PyFalse,
	"None": PyNil,
}

var PyBuiltInFuncMap = map[string]func(){}


// built in module

var ModuleGopy = Module{
	funcs: ModuleDict{
		"go": func(args *PyArgs) PyObject {
			return PyTrue
		},
	},
}


var ModuleTime = Module{
	name: "time",
	funcs: ModuleDict{
		"sleep": func(args *PyArgs) PyObject {
			time.Sleep(1 * time.Second)
			return PyTrue
		},
		"time": func(args *PyArgs) PyObject {
			return NewPyInt(time.Now().Unix())
		},
	},
}


var Modules = map[string]Module{
	"time": ModuleTime,
	"gopy": ModuleGopy,
}
