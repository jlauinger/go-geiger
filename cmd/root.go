package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/jlauinger/go-geiger/counter"
	"os"
)

var maxDepth int
var shortenSeenPackages, showStandardPackages, printLinkToPkgGoDev, printUnsafeLines, detailedStats, hideStats bool
var matchFilter, contextFilter string

var RootCmd = &cobra.Command{
	Use:   "geiger",
	Short: "Counts unsafe usages in dependencies",
	Long: `https://github.com/stg-tud/thesis-2020-lauinger-code/go-geiger`,
	Args: cobra.RangeArgs(0, 1000),
	Run: func(cmd *cobra.Command, args []string) {
		// run the go-geiger counter package Run function, which is the main entry point. Supply the configuration as
		// requested by the CLI parameters
		counter.Run(counter.Config{
			MaxDepth:             maxDepth,
			ShortenSeenPackages:  shortenSeenPackages,
			PrintLinkToPkgGoDev:  printLinkToPkgGoDev,

			DetailedStats:        detailedStats,
			HideStats:            hideStats,
			PrintUnsafeLines:     printUnsafeLines,

			ShowStandardPackages: showStandardPackages,
			MatchFilter:          matchFilter,
			ContextFilter:        contextFilter,

			Output:               os.Stdout,
		}, args...)
	},
}

func Execute() {
	// execute the root command with Cobra
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// register command line flags, for their meanings see their respective usage comments
	RootCmd.PersistentFlags().IntVarP(&maxDepth, "max-depth", "d", 10, "Maximum transitive import depth")
	RootCmd.PersistentFlags().BoolVar(&shortenSeenPackages, "show-only-once", true, "Do not repeat packages, show them only once and abbreviate further imports")
	RootCmd.PersistentFlags().BoolVarP(&printLinkToPkgGoDev, "link", "l",false, "Print link to pkg.go.dev instead of package name")

	RootCmd.PersistentFlags().BoolVarP(&detailedStats, "verbose", "v",false, "Show usage counts by different usage types")
	RootCmd.PersistentFlags().BoolVarP(&hideStats, "hide-stats", "q", false, "Hide statistics table, print only code. --show-code needs to be set manually")
	RootCmd.PersistentFlags().BoolVar(&printUnsafeLines, "show-code", false, "Print the code lines with unsafe usage")

	RootCmd.PersistentFlags().BoolVar(&showStandardPackages, "include-std", false, "Show / include Golang stdlib packages")
	RootCmd.PersistentFlags().StringVar(&matchFilter, "filter-match", "pointer", "Count only lines of requested match type (all,pointer,sizeof,offsetof,alignof,sliceheader,stringheader,uintptr). Default pointer")
	RootCmd.PersistentFlags().StringVar(&contextFilter, "filter-context", "all", "Count only lines of requested context type (all,variable,parameter,assignment,call,other). Default all")
}
