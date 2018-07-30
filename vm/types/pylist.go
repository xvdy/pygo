package types

import (
"fmt"
)

type PyList struct {
	PyTuple
}

func (pl *PyList) AsString() *string {
	str := fmt.Sprintf("[%s]", *pl.buildItemString())
	return &str
}

// Returns nil or exception object
func (pl *PyList) SetItem(key, value PyObject) PyObject {
	if _, ok := pl.attributes["__setitem__"]; ok { // set_func
		// found setattr, call it!
		panic("not implemented yet")
	}
	idxobj, isInt := key.(*PyInt)
	if !isInt {
		fmt.Sprintf("joo\n")
		return PyTypeError
	}

	if int(idxobj.value) > len(pl.Items) {
		return PyIndexError
	}

	pl.Items[idxobj.value] = value

	return nil
}

// Returns actual object or exception object
func (pl *PyList) GetItem(key PyObject) PyObject {
	if _, ok := pl.attributes["__getitem__"]; ok { // set_func
		// found setattr, call it!
		panic("not implemented yet")
	}
	idxobj, isInt := key.(*PyInt)
	if !isInt {
		fmt.Sprintf("joo\n")
		return PyTypeError
	}

	if int(idxobj.value) > len(pl.Items) {
		return PyIndexError
	}

	return pl.Items[idxobj.value]
}

func NewPyList(items []PyObject) PyObject {
	pl := &PyList{
		PyTuple{Items: items},
	}
	pl.PyObjInit()
	return pl
}

