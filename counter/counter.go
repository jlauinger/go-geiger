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

type Stats struct {
	ImportCount             int
	UnsafeCount             int
	TransitivelyUnsafeCount int
	SafeCount               int
}

type LocalPackageCounts struct {
	Local      int
	Variable   int
	Parameter  int
	Assignment int
	Call       int
	Other      int
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

	seenSelectors := map[*ast.SelectorExpr]bool{}

	inspectResult.WithStack([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.SelectorExpr)
		_, ok := seenSelectors[node]
		if ok {
			return true
		}
		seenSelectors[node] = true

		if !shouldCountSelectorExpr(node, config) {
			return true
		}

		if isInAssignment(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "assignment" {
				localPackageCounts.Local++
				localPackageCounts.Assignment++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isArgument(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "call" {
				localPackageCounts.Local++
				localPackageCounts.Call++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isParameter(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "parameter" {
				localPackageCounts.Local++
				localPackageCounts.Parameter++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isInVariableDefinition(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "variable" {
				localPackageCounts.Local++
				localPackageCounts.Variable++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if config.ContextFilter == "all" || config.ContextFilter == "other" {
			localPackageCounts.Local++
			localPackageCounts.Other++
			if config.PrintUnsafeLines {
				printLine(pkg, n)
			}
		}

		return true
	})

	seenIdents := map[*ast.Ident]bool{}

	inspectResult.WithStack([]ast.Node{(*ast.Ident)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.Ident)
		_, ok := seenIdents[node]
		if ok {
			return true
		}
		seenIdents[node] = true

		if !shouldCountIdent(node, config) {
			return true
		}

		if isInAssignment(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "assignment" {
				localPackageCounts.Local++
				localPackageCounts.Assignment++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isArgument(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "call" {
				localPackageCounts.Local++
				localPackageCounts.Call++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isParameter(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "parameter" {
				localPackageCounts.Local++
				localPackageCounts.Parameter++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isInVariableDefinition(stack) {
			if config.ContextFilter == "all" || config.ContextFilter == "variable" {
				localPackageCounts.Local++
				localPackageCounts.Variable++
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if config.ContextFilter == "all" || config.ContextFilter == "other" {
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

func shouldCountSelectorExpr(node *ast.SelectorExpr, config Config) bool {
	if (config.MatchFilter == "all" || config.MatchFilter == "pointer") && isUnsafePointer(node) {
		return true
	}
	if (config.MatchFilter == "all" || config.MatchFilter == "sizeof") && isUnsafeSizeof(node) {
		return true
	}
	if (config.MatchFilter == "all" || config.MatchFilter == "alignof") && isUnsafeAlignof(node) {
		return true
	}
	if (config.MatchFilter == "all" || config.MatchFilter == "offsetof") && isUnsafeOffsetof(node) {
		return true
	}
	if (config.MatchFilter == "all" || config.MatchFilter == "sliceheader") && isReflectSliceHeader(node) {
		return true
	}
	if (config.MatchFilter == "all" || config.MatchFilter == "stringheader") && isReflectStringHeader(node) {
		return true
	}
	return false
}

func shouldCountIdent(node *ast.Ident, config Config) bool {
	return (config.MatchFilter == "all" || config.MatchFilter == "uintptr") && isUintptr(node)
}