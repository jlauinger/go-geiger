package testdata

import "unsafe"

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

/*
 This package should have:

 total count: 5

 variable definition: 2
 parameter definition: 1
 assignment: 1
 call: 1
 */