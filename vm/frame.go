package vm

import (
	"pygo/vm/types"
	"fmt"
	"time"
	"io"
)

type PyFrame struct {
	stack    *PyObjStack
	blocks   *BlockStack
	names    map[string]types.PyObject // for fast locals!
	position int64
	funcs    map[string]types.PyFunc
}

func NewPyFrame(stacksize uint64) *PyFrame {
	return &PyFrame{
		stack:  NewPyObjStack(stacksize),
		blocks: NewBlockStack(10000),
		names:  make(map[string]types.PyObject),
	}
}

func (code *PyCode) eval(frame *PyFrame) (types.PyObject, error) {
	frame.position = 0 // Start from the beginning on eval!

	starttime := time.Now()
	defer func() {
		code.Log(fmt.Sprintf("Evaluation took %s.", time.Since(starttime)), true)
	}()

	for {
		// Get opcode
		op_position := frame.position
		code.code.setPos(frame.position)
		opcode, err := code.code.readByte()
		if err != nil {
			if err == io.EOF {
				code.Log("End of code reached", false)
				break
			} else {
				panic(err.Error())
			}
		}
		frame.position += 1
		code.vm.runtime.instructions += 1

		// Has opcode arguments?
		var oparg uint16
		if opcode >= types.HasArgLimes {
			oparg, err = code.code.readWord()
			if err != nil {
				panic("opcode requires arguments, but fetching failed: " + err.Error())
			}
			frame.position += 2
		}

		switch opcode {
		case types.POP_TOP:
			code.Log("Pop top", true)
			frame.stack.Pop()
		case types.ROT_TWO:
			code.Log("Rot two", true)
			op1stack := frame.stack.Pop()
			if op1stack == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}
			op2stack := frame.stack.Pop()
			if op2stack == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}
			frame.stack.Push(op1stack)
			frame.stack.Push(op2stack)
		case types.BINARY_MULTIPLY, types.INPLACE_MULTIPLY, types.INPLACE_ADD, types.BINARY_ADD, types.BINARY_SUBTRACT:
			code.Log(fmt.Sprintf("Beginning math operation (%d)", opcode), true)
			op1 := frame.stack.Pop()
			if op1 == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}

			op2 := frame.stack.Pop()
			if op2 == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}

			// save old content for log
			op2_old := *op2.AsString()

			var err types.PyObject
			var result types.PyObject

			switch opcode {
				// todo check type convert
				case types.INPLACE_MULTIPLY:
					err, result = op2.(*types.PyInt).Operation(types.OpMultiply, op1, true)
				case types.BINARY_MULTIPLY:
					err, result = op2.(*types.PyInt).Operation(types.OpMultiply, op1, false)
				case types.INPLACE_ADD:
					err, result = op2.(*types.PyInt).Operation(types.OpAdd, op1, true)
				case types.BINARY_ADD:
					err, result = op2.(*types.PyInt).Operation(types.OpAdd, op1, false)
				case types.BINARY_SUBTRACT:
					err, result = op2.(*types.PyInt).Operation(types.OpSubtract, op1, false)
				default:
					panic("Not implemented")
			}

			if err != nil {
				code.runtimeError("Exception raised: " + *err.(*types.PyException).Name) // TODO: handle correctly, this is only provisory
			}

			if result == nil {
				code.runtimeError("Result is nil in math operation")
			}

			code.Log(fmt.Sprintf("Operation finished: %v [%v] %v = %v", op2_old, opcode, *op1.AsString(), *result.AsString()), true)
			frame.stack.Push(result)
		case types.BINARY_SUBSCR:
			key := frame.stack.Pop()
			obj := frame.stack.Pop()
			if key == nil || obj == nil {
				code.runtimeError("Key or object is nil")
			}

			// todo map GetItem to different type such as PyTuple PyList
			result := obj.(*types.PyList).GetItem(key)
			frame.stack.Push(result)
		case types.STORE_SUBSCR:
			key := frame.stack.Pop()
			obj := frame.stack.Pop()
			value := frame.stack.Pop()
			fmt.Printf("%T\n", obj)
			err := obj.(*types.PyList).SetItem(key, value)
			if err != nil {
				code.runtimeError("Exception raised: " + *err.(*types.PyException).Name)
			}
		case types.PRINT_ITEM:
			code.Log("Print item", true)
			stackitem := frame.stack.Pop()
			if stackitem == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}
			//fmt.Printf("Item: %v\n", stackitem)
			fmt.Printf("%s ", *stackitem.AsString())
		case types.PRINT_NEWLINE:
			fmt.Println()
			code.Log("Print newline", true)
		case types.UNPACK_SEQUENCE:
			code.Log("Unpack sequence", true)
			stackitem := frame.stack.Pop()
			if stackitem == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}

			items := stackitem.GetValue().([]types.PyObject)
			for i := len(items) - 1; i >= 0; i-- {
				code.Log(fmt.Sprintf("Unpacking %d -> %v", i, items[i].GetValue()), true)
				frame.stack.Push(items[i])
			}
		case types.RETURN_VALUE:
			// TODO: Abschlussarbeiten durchführen, was zB? Siehe ceval.c in CPython
			stackitem := frame.stack.Pop()
			if stackitem == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}
			code.Log(fmt.Sprintf("Return value: %v", *stackitem.AsString()), true)
			return stackitem, nil
		case types.POP_BLOCK:
			code.Log("Pop block", true)
			frame.blocks.Pop()
		case types.LOAD_NAME:
			name := *code.names.(*types.PyTuple).Items[int(oparg)].GetValue().(*string)

			// Check wheter it's a built in type
			value, is_builtin := types.PyBuiltInTypeMap[name]
			if is_builtin {
				code.Log(fmt.Sprintf("Load built-in name %s (= %v, push to stack)", name, *value.AsString()), true)
				frame.stack.Push(value)
			} else {
				// It's no built-in, determine value

				_, global_found := code.vm.runtime.mainframe.names[name]

				// TODO CHECK: Not sure whether to give global or local priority
				if global_found {
					// Get global
				} else {

				}

				// Workaround: Get from local context
				value, local_found := frame.names[name]
				if !local_found {
					code.runtimeError(fmt.Sprintf("Could not find name (%v) in local namespace", name))
				}
				code.Log(fmt.Sprintf("Load name %s (= %v, push to stack)", name, *value.AsString()), true)
				frame.stack.Push(value)
			}
		case types.STORE_NAME:
			name := *code.names.(*types.PyTuple).Items[int(oparg)].GetValue().(*string)

			_, global_found := code.vm.runtime.mainframe.names[name]

			// TODO CHECK: Not sure whether to give global or local priority
			if global_found {
				// Set global
			} else {

			}

			// Workaround: Set in local context
			stackitem := frame.stack.Pop()
			if stackitem == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}
			frame.names[name] = stackitem

			code.Log(fmt.Sprintf("Store name: %s = %v", name, *frame.names[name].AsString()), true)
		case types.STORE_FAST:
			stackitem := frame.stack.Pop()
			if stackitem == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}
			name := *code.varnames.(*types.PyTuple).Items[int(oparg)].(*types.PyString).AsString()
			frame.names[name] = stackitem
			code.Log(fmt.Sprintf("Store FAST name: %s = %v", name, *frame.names[name].AsString()), true)
		case types.LOAD_FAST:
			name := *code.varnames.(*types.PyTuple).Items[int(oparg)].(*types.PyString).AsString()
			item, ok := frame.names[name]
			if !ok {
				code.runtimeError("Could not find item in varnames")
			}
			code.Log(fmt.Sprintf("Load FAST name: %s (= %v, pushing on stack)", name, *frame.names[name].AsString()), true)
			frame.stack.Push(item)
		case types.BUILD_TUPLE:
			items := make([]types.PyObject, oparg)
			for i := 0; i < int(oparg); i++ {
				stackitem := frame.stack.Pop()
				if stackitem == nil {
					code.runtimeError("Stackitem is nil during tuple build process")
				}
				items[i] = stackitem
			}
			fmt.Println("Creating TUPLE!")
			tuple := types.NewPyTuple(items)
			frame.stack.Push(tuple)
			code.Log(fmt.Sprintf("Build tuple (%d items: %s)", oparg, *tuple.(*types.PyTuple).AsString()), true)
		case types.BUILD_LIST:
			items := make([]types.PyObject, oparg)
			for i := 0; i < int(oparg); i++ {
				stackitem := frame.stack.Pop()
				if stackitem == nil {
					code.runtimeError("Stackitem is nil during list build process")
				}
				items[i] = stackitem
			}
			fmt.Println("Creating LIST!")
			list := types.NewPyList(items)
			frame.stack.Push(list)
			code.Log(fmt.Sprintf("Build list (%d items: %s)", oparg, *list.(*types.PyList).AsString()), true)
		case types.LOAD_ATTR:
			name := code.names.(*types.PyTuple).Items[int(oparg)].(*types.PyString)
			obj := frame.stack.Pop()

			result := obj.GetAttr(name, types.PyNil) // TODO: Check+raise Exception in result! This might return PyAttributeError

			if _, is_exception := result.(*types.PyException); is_exception {
				panic("Exception raised! To be handled correctly.")
			}

			frame.stack.Push(result)
			code.Log(fmt.Sprintf("Load attr [getattr(%v, %v) = %v]", *obj.AsString(), *name.AsString(), *result.AsString()), true)
		case types.COMPARE_OP:
			panic("not implement")
			//op1stack := frame.stack.Pop()
			//if op1stack == nil {
			//	code.runtimeError("Stackitem cannot be nil.")
			//}
			//op2stack := frame.stack.Pop()
			//if op2stack == nil {
			//	code.runtimeError("Stackitem cannot be nil.")
			//}
			//compareFunc, found := compareMap[int(oparg)]
			//if !found {
			//	code.runtimeError(fmt.Sprintf("Could not find compare function %d in map", oparg))
			//}
			//result := compareFunc(op2stack, op1stack)
			//code.Log(fmt.Sprintf("Compare op (%d); result = %t", oparg, result.IsTrue()), true)
			//frame.stack.Push(result)
		case types.LOAD_CONST:
			value := code.consts.(*types.PyTuple).Items[int(oparg)]
			code.Log(fmt.Sprintf("Load const: %v (pushing on stack)", *value.AsString()), true)
			frame.stack.Push(value)
		case types.LOAD_GLOBAL:
			name := *code.names.(*types.PyTuple).Items[int(oparg)].GetValue().(*string)

			// Check wheter it's a built in type
			value, is_builtin := types.PyBuiltInTypeMap[name]
			if is_builtin {
				code.Log(fmt.Sprintf("Load built-in name %s (= %v, push to stack)", name, *value.AsString()), true)
				frame.stack.Push(value)
			} else {
				// It's no built-in, determine value
				value, global_found := code.vm.runtime.mainframe.names[name]

				// TODO CHECK: Not sure whether to give global or local priority
				if !global_found {
					code.runtimeError(fmt.Sprintf("Could not find name (%v) in global namespace", name))
				}

				code.Log(fmt.Sprintf("Load GLOBAL name %s (= %v, push to stack)", name, *value.AsString()), true)
				frame.stack.Push(value)
			}
		case types.SETUP_LOOP:
			code.Log("Setup loop", true)
			frame.blocks.Push(op_position, frame.position+int64(oparg))
		case types.IMPORT_NAME:
			name := code.names.(*types.PyTuple).Items[int(oparg)].(*types.PyString).AsString()
			module := types.NewPyModule(name)
			frame.stack.Push(module)
			code.Log(fmt.Sprintf("Import name (%s = %v, pushed on stack)", *name, *module.AsString()), true)
		case types.JUMP_ABSOLUTE:
			code.Log(fmt.Sprintf("Jump absolute (%d)\n", oparg), true)
			frame.position = int64(oparg)
		case types.POP_JUMP_IF_FALSE:
			stackitem := frame.stack.Pop()
			if stackitem == nil {
				code.runtimeError("Stackitem cannot be nil.")
			}

			if !stackitem.IsTrue() {
				frame.position = int64(oparg)
			}

			code.Log(fmt.Sprintf("Pop/jump if false (result = %t)", stackitem.IsTrue()), true)
		case types.CALL_FUNCTION:
			code.Log(fmt.Sprintf("Call function (args=%d)", oparg), true)

			args := types.NewPyArgs()

			if oparg > 0 {
				nkwargs := (oparg >> 8) & 0xff
				// Keyword arguments first (high byte)
				for i := 0; i < int(nkwargs); i++ {
					arg := frame.stack.Pop()
					code.Log(fmt.Sprintf("   Received kw-arg: %v", arg), true)
					panic("Set kw! - TODO")
				}

				nposargs := oparg & 0xff
				// positional parameters (low byte)
				for i := 0; i < int(nposargs); i++ {
					arg := frame.stack.Pop()
					code.Log(fmt.Sprintf("   Received pos-arg: %v", arg), true)
					args.AddPositional(arg)
				}
			}

			fobj := frame.stack.Pop()
			if fobj == nil {
				code.runtimeError("Stack empty, expected: function object")
			}
			if fobj.(*types.PyFunc).IsExternal() {
				panic("not implement")
				//code.Log(fmt.Sprintf("Code obj argcount: %d", fobj.(*types.PyFunc).codeobj.(*PyCode).argcount), true)
			}
			result := fobj.(*types.PyFunc).Run(args) // TODO FIX!!
			frame.stack.Push(result)           // dunno?
		case types.MAKE_FUNCTION:
			code.Log(fmt.Sprintf("Make function (argcount=%d)", oparg), true)
			codeobj := frame.stack.Pop()
			fobj := NewPyFuncExternal(codeobj)
			if oparg > 0 {
				// Parse arguments
				panic("Not implemented yet")
				for i := 0; i < int(oparg); i++ {
					_ = frame.stack.Pop()
				}
			}
			frame.stack.Push(fobj)
		case types.MAKE_CLOSURE:
			panic("not implement")
			//code.Log(fmt.Sprintf("Make closure (argc = %d)", oparg), true)
			//c := frame.stack.Pop()
			//if c == nil {
			//	code.runtimeError("No code for closure (stackitem is nil)")
			//}
			//fn := types.NewPyFuncExternal(code)
			//closure := frame.stack.Pop()
			//if closure == nil {
			//	code.runtimeError("No closure for make closure (stackitem is nil)")
			//}
			//fn.(*types.PyFunc).SetClosure(closure)
			//
			//if oparg > 0 {
			//	// untested, but should work
			//	items := make([]types.PyObject, oparg)
			//	for i := 0; i < int(oparg); i++ {
			//		stackitem := frame.stack.Pop()
			//		if stackitem == nil {
			//			code.runtimeError("Stackitem is nil during tuple build process in make closure")
			//		}
			//		items = append(items, stackitem)
			//	}
			//	tuple := types.NewPyTuple(items)
			//	frame.stack.Push(tuple)
			//}
			//
			//frame.stack.Push(fn)
		case types.LOAD_CLOSURE:
			panic("not implement")
			//panic("Check for correct implementation")
			//obj := code.getFromCellFreeStorage(int(oparg))
			//frame.stack.Push(obj)
			//code.Log(fmt.Sprintf("Load closure (%s)", *obj.(*types.PyString).AsString()), true)
		case types.LOAD_DEREF:
			panic("Check for correct implementation")
			code.Log(fmt.Sprintf("Load deref (idx = %d)", oparg), true)
			obj := code.vm.runtime.freevars[int(oparg)]
			if obj == nil {
				code.runtimeError("Received freevar was nil")
			}
			frame.stack.Push(obj)
		case types.STORE_DEREF:
			panic("Check for correct implementation")
			code.Log(fmt.Sprintf("Store deref (idx = %d)", oparg), true)
			obj := frame.stack.Pop()
			code.vm.runtime.freevars[int(oparg)] = obj
		default:
			code.runtimeError(fmt.Sprintf("!!! Unhandled opcode: %c (%d)", opcode, opcode))
		}
	}
	//return PyNil, nil
	panic("unreachable")
}
