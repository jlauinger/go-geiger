package counter

import (
	"go/ast"
	"golang.org/x/tools/go/ast/inspector"
	"golang.org/x/tools/go/packages"
)

/**
 * this is a cache for package counts, so that they don't have to be recounted if packages are imported multiple times
 * through different paths
 */
var packageUnsafeCountCache map[*packages.Package]LocalPackageCounts

/*
 * initializes the cache to an empty hash map
 */
func initCache() {
	packageUnsafeCountCache = map[*packages.Package]LocalPackageCounts{}
}

/**
 * Stats represents the number of packages in the dependency tree of a given package, grouped by categories for directly
 * unsafe, transitively importing unsafe, and not using unsafe, as well as total
 */
type Stats struct {
	ImportCount             int
	UnsafeCount             int
	TransitivelyUnsafeCount int
	SafeCount               int
}

/**
 * LocalPackageCounts represents the number of unsafe usage sites in a given package, grouped by the context in which
 * they were found. The Local field represents the total number in the package, but it is not named Total in order to
 * disambiguate from the total number of unsafe usages in the package and its dependencies together.
 */
type LocalPackageCounts struct {
	Local      int
	Variable   int
	Parameter  int
	Assignment int
	Call       int
	Other      int
}

/**
 * analyzes a single package for unsafe usages and returns them using the LocalPackageCounts structure to group by
 * context
 */
func getUnsafeCount(pkg *packages.Package, config Config) LocalPackageCounts {
	// if this is a standard library package and they should be skipped, we can shortcut to returning a zero count value
	if config.ShowStandardPackages == false && isStandardPackage(pkg) {
		return LocalPackageCounts{}
	}

	// if the package was previously counted, we can shortcut by returning the counts from the cache
	cachedCounts, ok := packageUnsafeCountCache[pkg]
	if ok {
		return cachedCounts
	}

	// initialize an AST inspector object to make filtering nodes easier
	inspectResult := inspector.New(pkg.Syntax)
	localPackageCounts := LocalPackageCounts{}

	// this is a hash set preventing double-counting nodes, which can happen with the WithStack method
	seenSelectors := map[*ast.SelectorExpr]bool{}

	// then go over all selector expressions to find all unsafe and reflect match types. Use a stack for context
	// identification (which nodes are around an unsafe call site)
	inspectResult.WithStack([]ast.Node{(*ast.SelectorExpr)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		node := n.(*ast.SelectorExpr)
		// if this selector node was previously seen do not analyze it again. If not, register it as seen
		_, ok := seenSelectors[node]
		if ok {
			return true
		}
		seenSelectors[node] = true

		// check if this selector node is configured to be analyzed. If not continue with the next node
		if !shouldCountSelectorExpr(node, config) {
			return true
		}

		// check the context of this node
		if isInAssignment(stack) {
			// check if the assignment context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "assignment" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Assignment++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isArgument(stack) {
			// check if the argument context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "call" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Call++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isParameter(stack) {
			// check if the parameter context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "parameter" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Parameter++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isInVariableDefinition(stack) {
			// check if the variable context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "variable" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Variable++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if config.ContextFilter == "all" || config.ContextFilter == "other" {
			// otherwise, increment the Other category
			localPackageCounts.Local++
			localPackageCounts.Other++
			// if configured, also print the code line containing this node
			if config.PrintUnsafeLines {
				printLine(pkg, n)
			}
		}

		// return true to continue inspecting
		return true
	})

	// similarly to the selector expressions, this is a hash set of identifier nodes to avoid double-counting
	seenIdents := map[*ast.Ident]bool{}

	// go through all identifier nodes to catch uintptr nodes
	inspectResult.WithStack([]ast.Node{(*ast.Ident)(nil)}, func(n ast.Node, push bool, stack []ast.Node) bool {
		// if this identifier node was previously seen do not analyze it again. If not, register it as seen
		node := n.(*ast.Ident)
		_, ok := seenIdents[node]
		if ok {
			return true
		}
		seenIdents[node] = true

		// check if this identifier node is configured to be analyzed. If not continue with the next node
		if !shouldCountIdent(node, config) {
			return true
		}

		// check the context of this node
		if isInAssignment(stack) {
			// check if the assignment context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "assignment" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Assignment++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isArgument(stack) {
			// check if the argument context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "call" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Call++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isParameter(stack) {
			// check if the parameter context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "parameter" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Parameter++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if isInVariableDefinition(stack) {
			// check if the variable context is configured to be counted
			if config.ContextFilter == "all" || config.ContextFilter == "variable" {
				// if so, increment the corresponding counts
				localPackageCounts.Local++
				localPackageCounts.Variable++
				// if configured, also print the code line containing this node
				if config.PrintUnsafeLines {
					printLine(pkg, n)
				}
			}
		} else if config.ContextFilter == "all" || config.ContextFilter == "other" {
			// otherwise, increment the Other category
			localPackageCounts.Local++
			localPackageCounts.Other++
			// if configured, also print the code line containing this node
			if config.PrintUnsafeLines {
				printLine(pkg, n)
			}
		}

		// return true to continue inspecting
		return true
	})

	// store the counts in the cache and return them
	packageUnsafeCountCache[pkg] = localPackageCounts
	return localPackageCounts
}

/**
 * gets the total unsafe count for a package including its dependencies
 */
func getTotalUnsafeCount(pkg *packages.Package, config Config, seen *map[*packages.Package]bool) int {
	// if this package was already counted, return 0 to not count it again. This is imported so packages that get
	// imported multiple times do not let the unsafe count be higher than it should. If not, mark it as seen.
	_, ok := (*seen)[pkg]
	if ok {
		return 0
	}
	(*seen)[pkg] = true

	// start with the local package count
	totalCount := getUnsafeCount(pkg, config).Local

	// then go through all children and add their respective total count to the total count for this package
	for _, child := range pkg.Imports {
		totalCount += getTotalUnsafeCount(child, config, seen)
	}

	// finally return the resulting total count
	return totalCount
}

/**
 * returns true if a selector node should be counted as unsafe or reflect call site
 */
func shouldCountSelectorExpr(node *ast.SelectorExpr, config Config) bool {
	// return true only if it is a countable unsafe or reflect node and the respective node should be counted according
	// to the configuration
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
	// if none matched, return false to not count this node
	return false
}

/**
 * returns true if an identifier node should be counted as uintptr
 */
func shouldCountIdent(node *ast.Ident, config Config) bool {
	// return true only if uintptr instances should be counted by the configuration and this actually is a uintptr node
	return (config.MatchFilter == "all" || config.MatchFilter == "uintptr") && isUintptr(node)
}