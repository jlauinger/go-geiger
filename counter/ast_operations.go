package counter

import (
	"go/ast"
)

/**
  Types of unsafe usages to be counted

  - unsafe.Pointer
  - unsafe.Offsetof
  - unsafe.Sizeof
  - unsafe.Alignof
  - reflect.SliceHeader
  - reflect.StringHeader
  - uintptr

  - plain count
  - usage as variable definition, possibly in struct
  - usage as function parameter type
  - usage in assignment
  - usage call argument
*/

func isUnsafePointer(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Pointer" {
			return true
		}
	}
	return false
}

func isUnsafeSizeof(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Sizeof" {
			return true
		}
	}
	return false
}

func isUnsafeOffsetof(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Offsetof" {
			return true
		}
	}
	return false
}

func isUnsafeAlignof(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Alignof" {
			return true
		}
	}
	return false
}

func isReflectStringHeader(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "reflect" && node.Sel.Name == "StringHeader" {
			return true
		}
	}
	return false
}

func isReflectSliceHeader(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "reflect" && node.Sel.Name == "SliceHeader" {
			return true
		}
	}
	return false
}

func isUintptr(node *ast.Ident) bool {
	if node.Name == "uintptr" {
		return true
	}
	return false
}


func isArgument(stack []ast.Node) bool {
	// skip the last stack elements because the unsafe.Pointer SelectorExpr is itself a call expression.
	// the selector expression is in function position of a call, and we are not interested in that.
	for i := len(stack) - 2; i > 0; i-- {
		n := stack[i - 1]
		_, ok := n.(*ast.CallExpr)
		if ok {
			return true
		}
	}
	return false
}

func isInAssignment(stack []ast.Node) bool {
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		_, ok := n.(*ast.AssignStmt)
		if ok {
			return true
		}
		_, ok = n.(*ast.CompositeLit)
		if ok {
			return true
		}
		_, ok = n.(*ast.ReturnStmt)
		if ok {
			return true
		}
	}
	return false
}

func isParameter(stack []ast.Node) bool {
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		_, ok := n.(*ast.FuncType)
		if ok {
			return true
		}
	}
	return false
}

func isInVariableDefinition(stack []ast.Node) bool {
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		_, ok := n.(*ast.GenDecl)
		if ok {
			return true
		}
	}
	return false
}
