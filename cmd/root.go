package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/jlauinger/go-geiger/counter"
	"os"
)

var maxDepth int
var shortenSeenPackages, showStandardPackages, printLinkToPkgGoDev, printUnsafeLines, detailedStats, hideStats bool
var filter string

var RootCmd = &cobra.Command{
	Use:   "geiger",
	Short: "Counts unsafe usages in dependencies",
	Long: `https://github.com/stg-tud/thesis-2020-lauinger-code/go-geiger`,
	Args: cobra.RangeArgs(0, 1000),
	Run: func(cmd *cobra.Command, args []string) {
		counter.Run(counter.Config{
			MaxDepth:             maxDepth,
			ShortenSeenPackages:  shortenSeenPackages,
			PrintLinkToPkgGoDev:  printLinkToPkgGoDev,

			DetailedStats:        detailedStats,
			HideStats:            hideStats,
			PrintUnsafeLines:     printUnsafeLines,

			ShowStandardPackages: showStandardPackages,
			Filter:               filter,

			Output:               os.Stdout,
		}, args...)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&maxDepth, "max-depth", "d", 10, "Maximum transitive import depth")
	RootCmd.PersistentFlags().BoolVar(&shortenSeenPackages, "show-only-once", true, "Do not repeat packages, show them only once and abbreviate further imports")
	RootCmd.PersistentFlags().BoolVarP(&printLinkToPkgGoDev, "link", "l",false, "Print link to pkg.go.dev instead of package name")

	RootCmd.PersistentFlags().BoolVarP(&detailedStats, "verbose", "v",false, "Show usage counts by different usage types")
	RootCmd.PersistentFlags().BoolVarP(&hideStats, "hide-stats", "q", false, "Hide statistics table, print only code. --show-code needs to be set manually")
	RootCmd.PersistentFlags().BoolVar(&printUnsafeLines, "show-code", false, "Print the code lines with unsafe usage")

	RootCmd.PersistentFlags().BoolVar(&showStandardPackages, "include-std", false, "Show / include Golang stdlib packages")
	RootCmd.PersistentFlags().StringVarP(&filter, "filter", "f", "all", "Print only lines of requested type (variable,parameter,assignment,call,other). You need to specify --show-code also.")
}
