package types

import (
	"fmt"
)

type PyObject interface {
	GetValue() interface{}
	IsTrue() bool
	AsString() *string

	GetAttr(name, standard PyObject) PyObject

	//Operation(op int, obj2 PyObject, inplace bool) (PyObject, PyObject)

	// for list and tuple
	//SetItem(key, value PyObject) PyObject
	//GetItem(key PyObject) PyObject
}

type PyObjectData struct {
	attributes map[string]PyObject
}

func (obj *PyObjectData) PyObjInit() {
	obj.attributes = make(map[string]PyObject)
}

func (obj *PyObjectData) Operation(op int, obj2 PyObject, inplace bool) (PyObject, PyObject) {
	return PyTypeError, nil
}

func (obj *PyObjectData) GetAttr(name, standard PyObject) PyObject {
	name_string, ok := name.(*PyString)
	if !ok {
		panic(fmt.Sprintf("getattr(_, name [%v]) is no PyString", name))
	}
	value, found := obj.attributes[*name_string.AsString()]
	if !found {
		if standard != PyNil {
			return standard
		} else {
			return PyAttributeError
		}
	}
	return value
}

// Returns nil or exception object
func (obj *PyObjectData) SetItem(key, value PyObject) PyObject {
	if _, ok := obj.attributes["__setitem__"]; ok { // set_func
		// found setattr, call it!
		panic("not implemented yet")
	}
	panic("stop")
	return PyTypeError
}

// Returns actual object or exception object
func (obj *PyObjectData) GetItem(key PyObject) PyObject {
	if _, ok := obj.attributes["__getitem__"]; ok { // set_func
		// found setattr, call it!
		panic("not implemented yet")
	}
	return PyTypeError
}
