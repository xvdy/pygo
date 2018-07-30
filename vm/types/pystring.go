package types

type PyString struct {
	PyObjectData
	value *string
}

func (ps *PyString) GetValue() interface{} {
	return ps.value
}

func (ps *PyString) AsString() *string {
	return ps.value
}

func (ps *PyString) IsTrue() bool {
	return len(*ps.value) > 0
}


func NewPyString(value *string) PyObject {
	ps := &PyString{
		value: value,
	}
	ps.PyObjInit()
	return ps
}
