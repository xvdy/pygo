package types

import (
	"fmt"
	"strings"
)

type PyTuple struct {
	PyObjectData
	Items []PyObject
}

func (pt *PyTuple) AsString() *string {
	str := fmt.Sprintf("(%s)", *pt.buildItemString())
	return &str
}

func (pt *PyTuple) GetValue() interface{} {
	return pt.Items
}

func (pt *PyTuple) IsTrue() bool {
	return len(pt.Items) > 0
}

func (pt *PyTuple) buildItemString() *string {
	// TODO: Make it more performant
	str := strings.Join(func() []string {
		res := make([]string, MinInt(len(pt.Items), 11))
		for idx, item := range pt.Items {
			res[idx] = *item.AsString()
			if idx >= 9 && (idx+1) < (len(pt.Items)-1) {
				res[idx+1] = fmt.Sprintf("... %d more", len(pt.Items)-10)
				break
			}
		}
		return res
	}(), ", ")
	return &str
}

// Returns (PyException, resultobj)
func (pt *PyTuple) Operation(op int, obj2 PyObject, inplace bool) (PyObject, PyObject) {
	if !inplace {
		panic("Not implemented")
	}

	switch op {
	case OpMultiply:
		value, isInt := obj2.(*PyInt)
		if !isInt {
			fmt.Println("TypeError! Multiply on list/tuple can only be done with integers.")
			return PyTypeError, nil
		}

		newList := make([]PyObject, 0, len(pt.Items)*int(value.value))
		for i := 0; i < int(value.value); i++ {
			//fmt.Printf("i = %d\n", i)
			for _, item := range pt.Items {
				//fmt.Printf("%d = %v\n", i*(idx+1), *item.asString())
				newList = append(newList, item)
			}
		}
		pt.Items = newList

		//fmt.Printf("New size: %d, requested = %d\n", pt.length(), int(value.value) * oldSize)

		return nil, pt
	}
	return PyTypeError, nil
}

func NewPyTuple(items []PyObject) PyObject {
	pt := &PyTuple{
		Items: items,
	}
	pt.PyObjInit()
	return pt
}

func MinInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
