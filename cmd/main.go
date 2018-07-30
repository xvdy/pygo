package main

import (
	"fmt"
	"pygo/vm"
)

func main() {
	vm, err := vm.NewVM("/Users/xuetang/workspace/byterun/byterun/a.pyc", true)

	if err != nil {
		fmt.Printf("err:", err)
	}

	err = vm.Run()
	if err != nil {
		fmt.Printf("err:", err)
	}

}
