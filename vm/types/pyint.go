package types

import (
	"fmt"
	"strconv"
)

type PyInt struct {
	PyObjectData
	value int64
}

func (pi *PyInt) GetValue() interface{} {
	return pi.value
}

func (pi *PyInt) IsTrue() bool {
	return pi.value == 1
}

func (pi *PyInt) AsString() *string {
	s := strconv.Itoa(int(pi.value))
	return &s
}

// Returns (PyException, resulting object)
func (pi *PyInt) Operation(op int, obj2 PyObject, inplace bool) (PyObject, PyObject) {
	switch op {
	case OpMultiply:
		value, isInt := obj2.(*PyInt)
		if !isInt {
			fmt.Println("TypeError! Multiply on int can only be done with integers.")
			return PyTypeError, nil
		}
		if inplace {
			pi.value *= value.value
		} else {
			return nil, NewPyInt(pi.value * value.value)
		}
	case OpAdd:
		value, isInt := obj2.(*PyInt)
		if !isInt {
			fmt.Println("TypeError! Multiply on int can only be done with integers.")
			return PyTypeError, nil
		}
		if inplace {
			pi.value += value.value
		} else {
			return nil, NewPyInt(pi.value + value.value)
		}
	case OpSubtract:
		value, isInt := obj2.(*PyInt)
		if !isInt {
			fmt.Println("TypeError! Multiply on int can only be done with integers.")
			return PyTypeError, nil
		}

		if inplace {
			pi.value -= value.value
		} else {
			return nil, NewPyInt(pi.value - value.value)
		}
	default:
		return PyTypeError, nil
	}
	return nil, pi
}

func NewPyInt(value int64) PyObject {
	pi := &PyInt{
		value: value,
	}
	pi.PyObjInit()
	return pi
}
