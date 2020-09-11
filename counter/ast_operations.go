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

/**
 * returns true if the specified node is an unsafe.Pointer selector node by checking if the selector is an identifier
 * and the identifier names match.
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

/**
 * returns true if the specified node is an unsafe.Sizeof selector node by checking if the selector is an identifier
 * and the identifier names match.
 */
func isUnsafeSizeof(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Sizeof" {
			return true
		}
	}
	return false
}

/**
 * returns true if the specified node is an unsafe.Offsetof selector node by checking if the selector is an identifier
 * and the identifier names match.
 */
func isUnsafeOffsetof(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Offsetof" {
			return true
		}
	}
	return false
}

/**
 * returns true if the specified node is an unsafe.Alignof selector node by checking if the selector is an identifier
 * and the identifier names match.
 */
func isUnsafeAlignof(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "unsafe" && node.Sel.Name == "Alignof" {
			return true
		}
	}
	return false
}

/**
 * returns true if the specified node is a reflect.StringHeader selector node by checking if the selector is an identifier
 * and the identifier names match.
 */
func isReflectStringHeader(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "reflect" && node.Sel.Name == "StringHeader" {
			return true
		}
	}
	return false
}

/**
 * returns true if the specified node is a reflect.SliceHeader selector node by checking if the selector is an identifier
 * and the identifier names match.
 */
func isReflectSliceHeader(node *ast.SelectorExpr) bool {
	switch X := node.X.(type) {
	case *ast.Ident:
		if X.Name == "reflect" && node.Sel.Name == "SliceHeader" {
			return true
		}
	}
	return false
}

/**
 * returns true if the specified node is a uintptr identifier node.
 */
func isUintptr(node *ast.Ident) bool {
	return node.Name == "uintptr"
}


/**
 * returns true if the node referenced by the given node stack is part of a function call argument
 */
func isArgument(stack []ast.Node) bool {
	// skip the last stack elements because the unsafe.Pointer SelectorExpr is itself a call expression.
	// the selector expression is in function position of a call, and we are not interested in that.
	for i := len(stack) - 2; i > 0; i-- {
		n := stack[i - 1]
		// if there is a CallExpr node in the stack, we are in a function argument. Return true
		_, ok := n.(*ast.CallExpr)
		if ok {
			return true
		}
	}
	// otherwise, if no CallExpr could be found, return false
	return false
}

/**
 * returns true if the node referenced by the given node stack is part of an assignment
 */
func isInAssignment(stack []ast.Node) bool {
	// go up the node stack
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		// check if this node is an AssignStmt, CompositeLit or ReturnStmt, which count to being an assignment
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
	// if no node in the stack matched, this is not an assignment thus return false
	return false
}

/**
 * returns true if the node referenced by the given node stack is a parameter in a function definition
 */
func isParameter(stack []ast.Node) bool {
	// go up the node stack
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		// check if this node is a function declaration
		_, ok := n.(*ast.FuncType)
		if ok {
			return true
		}
	}
	// otherwise, if no node in the stack was a function declaration, this is not a parameter
	return false
}

/**
 * returns true if the node referenced by the given node stack is contained in a variable definition
 */
func isInVariableDefinition(stack []ast.Node) bool {
	// go up the node stack
	for i := len(stack); i > 0; i-- {
		n := stack[i - 1]
		// check if it is node of type generic declaration, which represents variable definitions
		_, ok := n.(*ast.GenDecl)
		if ok {
			return true
		}
	}
	// otherwise, if no node in the stack is of type GenDecl, this is not a variable definition
	return false
}
