package types

const (
	POP_TOP           = 1
	ROT_TWO           = 2
	BINARY_MULTIPLY   = 20
	BINARY_ADD        = 23
	BINARY_SUBTRACT   = 24
	BINARY_SUBSCR     = 25
	INPLACE_ADD       = 55
	INPLACE_MULTIPLY  = 57
	STORE_SUBSCR      = 60
	PRINT_ITEM        = 71
	PRINT_NEWLINE     = 72
	RETURN_VALUE      = 83
	POP_BLOCK         = 87
	STORE_NAME        = 90
	DELETE_NAME       = 91
	UNPACK_SEQUENCE   = 92
	FOR_ITER          = 93
	LOAD_CONST        = 100
	LOAD_NAME         = 101
	BUILD_TUPLE       = 102
	BUILD_LIST        = 103
	LOAD_ATTR         = 106
	COMPARE_OP        = 107
	IMPORT_NAME       = 108
	JUMP_ABSOLUTE     = 113
	POP_JUMP_IF_FALSE = 114
	LOAD_GLOBAL       = 116
	SETUP_LOOP        = 120
	LOAD_FAST         = 124
	STORE_FAST        = 125
	CALL_FUNCTION     = 131
	MAKE_FUNCTION     = 132
	MAKE_CLOSURE      = 134
	LOAD_CLOSURE      = 135
	LOAD_DEREF        = 136
	STORE_DEREF       = 137
)

const HasArgLimes = 90

const (
	OpAdd      = iota
	OpMultiply
	OpSubtract
)
