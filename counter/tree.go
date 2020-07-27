package counter

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/go/packages"
	"sort"
	"strconv"
)

type IndentType int
const (
	Space IndentType = iota
	I
	T
	L
)

func printPkgTree(pkg *packages.Package, indents []IndentType, config Config, table *tablewriter.Table,
	seen *map[*packages.Package]bool) (stats Stats) {
	(*seen)[pkg] = true

	countsInThisPackage := getUnsafeCount(pkg, config)
	totalCount := getTotalUnsafeCount(pkg, config, &map[*packages.Package]bool{})
	nameString := fmt.Sprintf("%s%s", getIndentString(indents), getPrintedPackageName(pkg, config))

	colors := getColors(countsInThisPackage.Local, totalCount, config)

	if config.DetailedStats && config.ContextFilter == "all" {
		table.Rich([]string{strconv.Itoa(totalCount), strconv.Itoa(countsInThisPackage.Local),
			strconv.Itoa(countsInThisPackage.Variable), strconv.Itoa(countsInThisPackage.Parameter),
			strconv.Itoa(countsInThisPackage.Assignment), strconv.Itoa(countsInThisPackage.Call),
			strconv.Itoa(countsInThisPackage.Other),
			nameString}, colors)
	} else {
		table.Rich([]string{strconv.Itoa(totalCount), strconv.Itoa(countsInThisPackage.Local), nameString}, colors)
	}

	childCount, _ := getImportsCount(pkg.Imports, config)
	nextIndents := getNextIndents(indents)

	if len(indents) == config.MaxDepth && childCount > 0 {
		if config.DetailedStats && config.ContextFilter == "all" {
			table.Append([]string{"", "", "", "", "", "", "", fmt.Sprintf("%sMaximum depth reached. Use --max-depth= to increase it",
				getIndentString(append(nextIndents, L)))})
		} else {
			table.Append([]string{"", "", fmt.Sprintf("%sMaximum depth reached. Use --max-depth= to increase it",
				getIndentString(append(nextIndents, L)))})
		}
		return
	}

	childKeys := make([]string, 0, len(pkg.Imports))
	for childKey := range pkg.Imports {
		childKeys = append(childKeys, childKey)
	}
	sort.Strings(childKeys)

	// do not count the outermost parent package in the import stats
	stats.ImportCount += childCount
	if countsInThisPackage.Local > 0 {
		stats.UnsafeCount += 1
	} else if totalCount > 0 {
		stats.TransitivelyUnsafeCount += 1
	} else {
		stats.SafeCount += 1
	}

	childIndex := 0
	for _, childKey := range childKeys {
		child := pkg.Imports[childKey]

		if config.ShowStandardPackages == false && isStandardPackage(child) {
			continue
		}

		childIndex++
		childIndents := getChildIndents(childIndex, childCount, nextIndents)

		_, ok := (*seen)[child]
		if config.ShortenSeenPackages && ok {
			countsInChild := getUnsafeCount(child, config)
			totalCountInChild := getTotalUnsafeCount(child, config, &map[*packages.Package]bool{})

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

			stats.ImportCount -= 1
			continue
		}

		childStats := printPkgTree(child, childIndents, config, table, seen)

		stats.ImportCount += childStats.ImportCount
		stats.UnsafeCount += childStats.UnsafeCount
		stats.TransitivelyUnsafeCount += childStats.TransitivelyUnsafeCount
		stats.SafeCount += childStats.SafeCount
	}

	return
}

func getChildIndents(childIndex int, childCount int, nextIndents []IndentType) []IndentType {
	isLast := childIndex == childCount

	var nextChildIndents []IndentType
	if isLast {
		nextChildIndents = append(nextIndents, L)
	} else {
		nextChildIndents = append(nextIndents, T)
	}
	return nextChildIndents
}

func getColors(countInThisPackage int, totalCount int, config Config) []tablewriter.Colors {
	var color int
	if countInThisPackage > 0 {
		color = tablewriter.FgRedColor
	} else if totalCount == 0 {
		color = tablewriter.FgGreenColor
	} else {
		color = tablewriter.Normal
	}
	if config.DetailedStats && config.ContextFilter == "all" {
		return []tablewriter.Colors{{color}, {color}, {color}, {color}, {color}, {color}, {color}, {color}}
	} else {
		return []tablewriter.Colors{{color}, {color}, {color}}
	}
}

func getNextIndents(indents []IndentType) []IndentType {
	var nextIndents []IndentType
	if len(indents) > 0 {
		nextIndents = indents[0 : len(indents)-1]
		if indents[len(indents)-1] == L || indents[len(indents)-1] == Space {
			nextIndents = append(nextIndents, Space)
		} else {
			nextIndents = append(nextIndents, I)
		}
	} else {
		nextIndents = []IndentType{}
	}
	return nextIndents
}

func getIndentString(indents []IndentType) string {
	str := ""
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

