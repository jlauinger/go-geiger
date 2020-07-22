package counter

import (
	"go/ast"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

var packageUnsafeCountCache map[*packages.Package]LocalPackageCounts

func initCache() {
	packageUnsafeCountCache = map[*packages.Package]LocalPackageCounts{}
}

/**
  Types of unsafe usages to be counted

  - plain unsafe.Pointer count

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

func getUnsafeCount(pkg *packages.Package, config Config) LocalPackageCounts {
	if config.ShowStandardPackages == false && isStandardPackage(pkg) {
		return LocalPackageCounts{}
	}

	cachedCounts, ok := packageUnsafeCountCache[pkg]
	if ok {
		return cachedCounts
	}

	inspectResult := inspector.New(pkg.Syntax)
	localPackageCounts := LocalPackageCounts{}

	seen := map[*ast.SelectorExpr]bool{}

	inspectResult.WithStack([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.SelectorExpr)
		_, ok := seen[node]
		if ok {
			return true
		}
		seen[node] = true

		if !isUnsafePointer(node) {
			return true
		}

		if isInAssignment(stack) {
			if config.Filter == "all" || config.Filter == "assignment" {
				localPackageCounts.Local++
				localPackageCounts.Assignment++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isArgument(stack) {
			if config.Filter == "all" || config.Filter == "call" {
				localPackageCounts.Local++
				localPackageCounts.Call++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isParameter(stack) {
			if config.Filter == "all" || config.Filter == "parameter" {
				localPackageCounts.Local++
				localPackageCounts.Parameter++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isInVariableDefinition(stack) {
			if config.Filter == "all" || config.Filter == "variable" {
				localPackageCounts.Local++
				localPackageCounts.Variable++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if config.Filter == "all" || config.Filter == "other" {
			localPackageCounts.Local++
			localPackageCounts.Other++
			if config.PrintUnsafeLines {
				printLine(pkg, n)
			}
		}

		return true
	})

	packageUnsafeCountCache[pkg] = localPackageCounts

	return localPackageCounts
}

func getTotalUnsafeCount(pkg *packages.Package, config Config, seen *map[*packages.Package]bool) int {
	_, ok := (*seen)[pkg]
	if ok {
		return 0
	}
	(*seen)[pkg] = true

	totalCount := getUnsafeCount(pkg, config).Local

	for _, child := range pkg.Imports {
		totalCount += getTotalUnsafeCount(child, config, seen)
	}

	return totalCount
}