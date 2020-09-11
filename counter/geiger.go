package counter

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/tools/go/packages"
	"io"
	"strings"
)

/**
 * The Config structure represents the command line arguments passed to go-geiger
 */
type Config struct {
	MaxDepth             int			// maximum depth in the import tree
	ShortenSeenPackages  bool			// if true, packages that are duplicated will only print their subtree once
	PrintLinkToPkgGoDev  bool			// if true, print a link to pkg.go.dev instead of the package name

	DetailedStats        bool			// if true, output separate numbers for different usage contexts
	HideStats            bool			// if true, do not output the table with usage counts but only print code
	PrintUnsafeLines     bool			// if true, print code lines where unsafe usages are found

	ShowStandardPackages bool			// if true, also show packages from the Go standard library which are omitted otherwise
	MatchFilter          string			// filter unsafe matches to all,pointer,sizeof,alignof,offsetof,sliceheader,stringheader, or uintptr
	ContextFilter        string			// filter context matches to all,variable,parameter,assignment,call,other

	Output               io.Writer		// output stream, needed to redirect output in testing
}

/**
 * Run is the main entry point for the go-geiger logic
 */
func Run(config Config, paths... string) {
	// set up the parsing mode: we need almost all parsing, only Types can be left out if we should not print lines
	mode := packages.NeedImports | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedFiles | packages.NeedName
	if config.PrintUnsafeLines {
		mode |= packages.NeedTypes
	}

	// load and parse requested packages
	pkgs, err := packages.Load(&packages.Config{
		Mode:       mode,
		Tests:      false,
	}, paths...)
	if err != nil {
		panic(err)
	}

	// analyze each package on its own
	for _, pkg := range pkgs {
		// reset the cache that is used to avoid recounting packages that are imported through different paths
		initCache()

		// initialize a table output with columns according to the configuration options
		table := tablewriter.NewWriter(config.Output)
		if config.DetailedStats && config.ContextFilter == "all" {
			table.SetHeader([]string{"With Dependencies", "Local Package", "Variable", "Parameter", "Assignment", "Call", "Other", "Package Path"})
			table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER,
				tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER,
				tablewriter.ALIGN_LEFT})
		} else if config.ContextFilter == "all" {
			table.SetHeader([]string{"With Dependencies", "Local Package", "Package Path"})
			table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT})
		} else {
			table.SetHeader([]string{"With Dependencies", fmt.Sprintf("Local Package %s", config.ContextFilter), "Package Path"})
			table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_LEFT})
		}
		table.SetBorder(false)
		table.SetColumnSeparator(" ")
		table.SetAutoWrapText(false)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

		// analyze package and print its dependency tree
		stats := printPkgTree(pkg, []IndentType{}, config, table, &map[*packages.Package]bool{})

		// if code lines were printed and stats should follow, add a blank line to separate them
		if config.PrintUnsafeLines && !config.HideStats {
			_, _ = fmt.Fprintln(config.Output)
		}

		// print stats table and summary if configured to do so
		if !config.HideStats {
			table.Render()
			printStats(pkg, stats, config)
		}
	}

	// print the legend to the stats table if configured to do so
	if !config.HideStats {
		printMatchConfig(config)
		printLegend(config)
	}
}

/**
 * prints a readable version of the match config
 */
func printMatchConfig(config Config) {
	countedTypes := make([]string, 0)

	// append all matches that are active to a list
	if config.MatchFilter == "all" || config.MatchFilter == "pointer" {
		countedTypes = append(countedTypes, "unsafe.Pointer")
	}
	if config.MatchFilter == "all" || config.MatchFilter == "sizeof" {
		countedTypes = append(countedTypes, "unsafe.Sizeof")
	}
	if config.MatchFilter == "all" || config.MatchFilter == "offsetof" {
		countedTypes = append(countedTypes, "unsafe.Offsetof")
	}
	if config.MatchFilter == "all" || config.MatchFilter == "alignof" {
		countedTypes = append(countedTypes, "unsafe.Alignof")
	}
	if config.MatchFilter == "all" || config.MatchFilter == "sliceheader" {
		countedTypes = append(countedTypes, "reflect.SliceHeader")
	}
	if config.MatchFilter == "all" || config.MatchFilter == "stringheader" {
		countedTypes = append(countedTypes, "reflect.StringHeader")
	}
	if config.MatchFilter == "all" || config.MatchFilter == "uintptr" {
		countedTypes = append(countedTypes, "uintptr")
	}

	// then join the list with commas and print it
	_, _ = fmt.Fprintf(config.Output, "Couting occurances of %s\n\n", strings.Join(countedTypes, ","))
}

/**
 * prints a colored legend of the meaning of the colors in the output table
 */
func printLegend(config Config) {
	_, _ = fmt.Fprintf(config.Output, "%s have no unsafe usages\n", color.GreenString("Packages in green"))
	_, _ = fmt.Fprintf(config.Output, "%s contain unsafe usages\n", color.RedString("Packages in red"))
	_, _ = fmt.Fprintf(config.Output, "%s import packages with unsafe usages\n", color.WhiteString("Packages in white"))
}

/**
 * prints summary stats about the usages found in the package and its dependencies
 */
func printStats(pkg *packages.Package, stats Stats, config Config) {
	// add a blank line for better readability
	_, _ = fmt.Fprintln(config.Output)

	// print a summary of package count. Total packages are the imported package and the package itself
	_, _ = fmt.Fprintf(config.Output, "Package %s including imports effectively makes up %d packages\n", pkg.PkgPath, stats.ImportCount+1)

	// print how many of those directly contain unsafe usages, if any
	if stats.UnsafeCount > 0 {
		_, _ = fmt.Fprint(config.Output, color.RedString("  %d of those contain unsafe usages\n", stats.UnsafeCount))
	}
	// print how many contain unsafe usages in their transitive dependencies, if any
	if stats.TransitivelyUnsafeCount > 0 {
		_, _ = fmt.Fprint(config.Output, color.WhiteString("  %d of those further import packages that contain unsafe usages\n",
			stats.TransitivelyUnsafeCount))
	}
	// and print how many do not use unsafe at all, if any
	if stats.SafeCount > 0 {
		_, _ = fmt.Fprint(config.Output, color.GreenString("  %d of those do not contain any unsafe usages\n", stats.SafeCount))
	}

	// then end with another blank line
	_, _ = fmt.Fprintln(config.Output)
}
