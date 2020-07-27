package testdata

import (
	"reflect"
	"unsafe"
)

type Bar struct {
	Baz unsafe.Pointer
}

func Foo(X unsafe.Pointer) {
	_ = unsafe.Pointer(&X)
}

func Entry() {
	x := 42
	var y unsafe.Pointer
	Foo(unsafe.Pointer(&x))
	Foo(y)
}

func SecondFunction() {
	x := 42
	_ = unsafe.Offsetof(x)
	_ = unsafe.Alignof(x)
	_ = unsafe.Sizeof(x)
	_ = (*reflect.SliceHeader)(unsafe.Pointer(&x))
	_ = (*reflect.StringHeader)(unsafe.Pointer(&x))
	var foo uintptr
	_ = foo
}

/*
 With matchType="pointer", this package should have:

 total count: 7

 variable definition: 2
 parameter definition: 1
 assignment: 3
 call: 1

 With matchType="all", this package should have:

 total count: 13

 variable definition: 3
 parameter definition: 1
 assignment: 8
 call: 1
 */