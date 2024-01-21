package iface

import "golang.org/x/exp/constraints"

type SignedSlice interface {
	~[]int | ~[]int8 | ~[]int16 | ~[]int32 | ~[]int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// If future releases of Go add new predeclared unsigned integer types,
// this constraint will be modified to include them.
type UnsignedSlice interface {
	~[]uint | ~[]uint8 | ~[]uint16 | ~[]uint32 | ~[]uint64 | ~[]uintptr
}

// Integer is a constraint that permits any integer type.
// If future releases of Go add new predeclared integer types,
// this constraint will be modified to include them.
type IntegerSlice interface {
	SignedSlice | UnsignedSlice
}

// Float is a constraint that permits any floating-point type.
// If future releases of Go add new predeclared floating-point types,
// this constraint will be modified to include them.
type FloatSlice interface {
	~[]float32 | ~[]float64
}

// Complex is a constraint that permits any complex numeric type.
// If future releases of Go add new predeclared complex numeric types,
// this constraint will be modified to include them.
type ComplexSlice interface {
	~[]complex64 | ~[]complex128
}

// Ordered is a constraint that permits any ordered type: any type
// that supports the operators < <= >= >.
// If future releases of Go add new ordered types,
// this constraint will be modified to include them.
type OrderedSlice interface {
	IntegerSlice | FloatSlice | ~[]string
}

type TypeWhere interface {
	constraints.Ordered | OrderedSlice | []interface{} | interface{}
}
