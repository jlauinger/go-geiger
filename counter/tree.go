package counter

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/go/packages"
	"sort"
	"strconv"
)

// this enum represents different possible indents used to build an ASCII tree shape
type IndentType int
const (
	Space IndentType = iota
	I
	T
	L
)

/**
 * this is the main function to build and print the dependency tree with counts for a given package
 */
func printPkgTree(pkg *packages.Package, indents []IndentType, config Config, table *tablewriter.Table,
	seen *map[*packages.Package]bool) (stats Stats) {

	// flag this package as seen in the hash set to prevent infinite loops
	(*seen)[pkg] = true

	// count this package for unsafe usages
	countsInThisPackage := getUnsafeCount(pkg, config)
	totalCount := getTotalUnsafeCount(pkg, config, &map[*packages.Package]bool{})
	// prepare the printed package name by adding the correct indent to form a tree and the package name
	nameString := fmt.Sprintf("%s%s", getIndentString(indents), getPrintedPackageName(pkg, config))

	// determine the color this package should get based on the counts for this package. The color is expressed as a
	// list because the number of columns in the table depends on the configuration
	colors := getColors(countsInThisPackage.Local, totalCount, config)

	// build the table row for this package with columns depending on the configuration
	if config.DetailedStats && config.ContextFilter == "all" {
		table.Rich([]string{strconv.Itoa(totalCount), strconv.Itoa(countsInThisPackage.Local),
			strconv.Itoa(countsInThisPackage.Variable), strconv.Itoa(countsInThisPackage.Parameter),
			strconv.Itoa(countsInThisPackage.Assignment), strconv.Itoa(countsInThisPackage.Call),
			strconv.Itoa(countsInThisPackage.Other),
			nameString}, colors)
	} else {
		table.Rich([]string{strconv.Itoa(totalCount), strconv.Itoa(countsInThisPackage.Local), nameString}, colors)
	}

	// determine the number of children to analyze, this depends on whether standard library packages should be included
	childCount, _ := getImportsCount(pkg.Imports, config)
	// and get the indents for children of this package, which are based on the position of the package in the dependency
	// ASCII tree
	nextIndents := getNextIndents(indents)

	// The indents also function as a convenient maximum depth condition. Check if we reached the maximum but there are
	// more children that are omitted
	if len(indents) == config.MaxDepth && childCount > 0 {
		// append a notice of missing packages to the output table. Depending on the configuration, the number of empty
		// columns in this row varies
		if config.DetailedStats && config.ContextFilter == "all" {
			table.Append([]string{"", "", "", "", "", "", "", fmt.Sprintf("%sMaximum depth reached. Use --max-depth= to increase it",
				getIndentString(append(nextIndents, L)))})
		} else {
			table.Append([]string{"", "", fmt.Sprintf("%sMaximum depth reached. Use --max-depth= to increase it",
				getIndentString(append(nextIndents, L)))})
		}
		// then do not continue analyzing any more children
		return
	}

	// get the package paths of all imported children and sort them alphabetically to achieve a consistent output
	childKeys := make([]string, 0, len(pkg.Imports))
	for childKey := range pkg.Imports {
		childKeys = append(childKeys, childKey)
	}
	sort.Strings(childKeys)

	// do not count the outermost parent package in the import stats. Otherwise, adjust the number of packages in the
	// respective categories by adding this package
	stats.ImportCount += childCount
	if countsInThisPackage.Local > 0 {
		stats.UnsafeCount += 1
	} else if totalCount > 0 {
		stats.TransitivelyUnsafeCount += 1
	} else {
		stats.SafeCount += 1
	}

	// then go over the imported children packages
	childIndex := 0
	for _, childKey := range childKeys {
		child := pkg.Imports[childKey]

		// if this is a standard library package and it should not explicitly be included, skip it
		if config.ShowStandardPackages == false && isStandardPackage(child) {
			continue
		}

		// get the correct indent for this child. This depends on whether it is the last in the tree
		childIndex++
		childIndents := getChildIndents(childIndex, childCount, nextIndents)

		// check if this package was already analyzed previously, which is possible because the same package can be
		// imported through different paths. In this case, we could shortcut to not printing it again if configured
		// to do so, which is currently the default
		_, ok := (*seen)[child]
		if config.ShortenSeenPackages && ok {
			// get the unsafe counts for this child, which will come from the cache
			countsInChild := getUnsafeCount(child, config)
			totalCountInChild := getTotalUnsafeCount(child, config, &map[*packages.Package]bool{})

			// build the table row for this child with columns depending on the configuration
			if config.DetailedStats && config.ContextFilter == "all" {
				table.Rich([]string{strconv.Itoa(totalCountInChild), strconv.Itoa(countsInChild.Local),
					strconv.Itoa(countsInChild.Variable), strconv.Itoa(countsInChild.Parameter),
					strconv.Itoa(countsInChild.Assignment), strconv.Itoa(countsInChild.Call),
					strconv.Itoa(countsInChild.Other),
					fmt.Sprintf("%s%s...", getIndentString(childIndents), getPrintedPackageName(child, config))},
					getColors(0, totalCountInChild, config))
			} else {
				table.Rich([]string{strconv.Itoa(totalCountInChild), strconv.Itoa(countsInChild.Local),
					fmt.Sprintf("%s%s...", getIndentString(childIndents), getPrintedPackageName(child, config))},
					getColors(0, totalCountInChild, config))
			}

			// reduce the import count because this package was already imported through a different path before and
			// then skip the recounting of the child
			stats.ImportCount -= 1
			continue
		}

		// if the child was not seen before, recursively analyze it
		childStats := printPkgTree(child, childIndents, config, table, seen)

		// add the numbers of packages by category of this child to the parent
		stats.ImportCount += childStats.ImportCount
		stats.UnsafeCount += childStats.UnsafeCount
		stats.TransitivelyUnsafeCount += childStats.TransitivelyUnsafeCount
		stats.SafeCount += childStats.SafeCount
	}

	return
}

/**
 * determines the correct indent sequence for a child package to form a visual ASCII tree
 */
func getChildIndents(childIndex int, childCount int, nextIndents []IndentType) []IndentType {
	// check if the package is the last of this layer in the tree
	isLast := childIndex == childCount
	var nextChildIndents []IndentType
	if isLast {
		// if so, append an L shape
		nextChildIndents = append(nextIndents, L)
	} else {
		// otherwise, since there will be more rows after, append a T shape
		nextChildIndents = append(nextIndents, T)
	}
	return nextChildIndents
}

/**
 * gets the colors for each column of a table row based on the unsafe counts
 */
func getColors(countInThisPackage int, totalCount int, config Config) []tablewriter.Colors {
	var color int
	// determine the row color by checking if the package has direct unsafe usages, transitive usages, or none
	if countInThisPackage > 0 {
		color = tablewriter.FgRedColor
	} else if totalCount == 0 {
		color = tablewriter.FgGreenColor
	} else {
		color = tablewriter.Normal
	}
	// depending on the configuration, return a list with appropriately many columns in this color
	if config.DetailedStats && config.ContextFilter == "all" {
		return []tablewriter.Colors{{color}, {color}, {color}, {color}, {color}, {color}, {color}, {color}}
	} else {
		return []tablewriter.Colors{{color}, {color}, {color}}
	}
}

/**
 * builds the indent string prefix for children at a given position, taking care of the current ASCII tree shape
 */
func getNextIndents(indents []IndentType) []IndentType {
	var nextIndents []IndentType
	// check if there already are any indents
	if len(indents) > 0 {
		// if so, take all but the last current indent symbol
		nextIndents = indents[0 : len(indents)-1]
		// then check if the last symbol was an L shape or a space
		if indents[len(indents)-1] == L || indents[len(indents)-1] == Space {
			// under an L or space should follow nothing, so add a space symbol
			nextIndents = append(nextIndents, Space)
		} else {
			// under I and T shapes should follow an I shape to carry on the connection for following rows
			nextIndents = append(nextIndents, I)
		}
	} else {
		// otherwise, initialize an empty list
		nextIndents = []IndentType{}
	}
	return nextIndents
}

/**
 * turns the list of indent tokens into an actual indent string using ASCII shapes
 */
func getIndentString(indents []IndentType) string {
	str := ""
	// go through the indents and append a suitable string representation of the tree shape symbol
	for _, indent := range indents {
		switch indent {
		case Space:
			str = fmt.Sprintf("%s%s", str, "  ")
		case I:
			str = fmt.Sprintf("%s%s", str, "│ ")
		case T:
			str = fmt.Sprintf("%s%s", str, "├─")
		case L:
			str = fmt.Sprintf("%s%s", str, "└─")
		}
	}
	return str
}

