package types

type PyNone struct {
	PyObjectData
}

func (pn *PyNone) GetValue() interface{} {
	return nil
}

func (pn *PyNone) IsTrue() bool {
	return false
}

func (pn *PyNone) AsString() *string {
	str := "None"
	return &str
}

func NewPyNone() PyObject {
	pn := new(PyNone)
	pn.PyObjInit()
	return pn
}
