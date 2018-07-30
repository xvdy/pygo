package types


var PyAttributeError = NewPyException("AttributeError")
var PyNameError = NewPyException("NameError")
var PyTypeError = NewPyException("TypeError")
var PyIndexError = NewPyException("IndexError")

type PyException struct {
	PyObjectData
	Name *string
	Msg  *string
}

func (pe *PyException) AsString() *string {
	return pe.Name
}

func (pe *PyException) GetValue() interface{} {
	return pe
}

func (pe *PyException) IsTrue() bool {
	return true
}

func NewPyException(name string) PyObject {
	excp := new(PyException)
	excp.PyObjInit()
	excp.Name = &name
	return excp
}

