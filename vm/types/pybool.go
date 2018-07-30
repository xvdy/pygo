package types

type PyBool struct {
	PyObjectData
	value bool
}

func (pb *PyBool) GetValue() interface{} {
	return pb.value
}

func (pb *PyBool) IsTrue() bool {
	return pb.value == true
}

func (pb *PyBool) AsString() *string {
	var str string
	if pb.value {
		str = "True"
	} else {
		str = "False"
	}
	return &str
}

func NewPyBool(value bool) PyObject {
	pb := &PyBool{
		value: value,
	}
	pb.PyObjInit()
	return pb
}
