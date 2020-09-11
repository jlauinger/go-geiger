package counter

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"strings"
)

/**
 * gets the package name or its link to pkg.go.dev, depending on the configuration
 */
func getPrintedPackageName(pkg *packages.Package, config Config) string {
	// check if links are requested and the current package looks like a package that can be associated with an URL
	// (i.e. that is not local to the current app). This is only a heuristic but it works reasonably well.
	if config.PrintLinkToPkgGoDev && pathLooksLikeAUrl(pkg.PkgPath) {
		return fmt.Sprintf("https://pkg.go.dev/%s", pkg.PkgPath)
	} else {
		// otherwise, return the plain package import path
		return pkg.PkgPath
	}
}

/**
 * checks whether a package path looks like a URL and can therefore be turned into a pkg.go.dev link
 */
func pathLooksLikeAUrl(path string) bool {
	// split the path by slashes
	components := strings.Split(path, "/")

	// if there are no slashes, it cannot be a URL-style pcakage name
	if len(components) <= 1 {
		return false
	}

	// otherwise, take the first component which should be the domain (e.g. github.com)
	domain := components[0]

	// finally, check if the domain contains a dot, which is a heuristic for a public registry instead of a local
	// import from the current app
	return strings.Contains(domain, ".")
}

/**
 * gets the number of imported packages but taking care of only counting packages that are requested to be analyzed
 */
func getImportsCount(pkgs map[string]*packages.Package, config Config) (childCount, stdLibCount int) {
	// iterate over all the packages
	for _, pkg := range pkgs {
		// check if it is part of the Go standard library
		if isStandardPackage(pkg) {
			// if so, it needs to be counted as a stdlib package
			stdLibCount++
			// but only as a child too if standard packages are configured to be part of the analysis
			if config.ShowStandardPackages {
				childCount++
			}
		} else {
			// otherwise, it only needs to be counted as a child
			childCount++
		}
	}
	return
}

/**
 * prints the line of code containing a specific AST node
 */
func printLine(pkg *packages.Package, n ast.Node) {
	// get the file and line number containing the AST node from the parsed package structure
	file := pkg.Fset.File(n.Pos())
	lineNumber := file.Position(n.Pos()).Line  // 1-based

	// determine the start and end offset in bytes that make up the line containing the AST node
	start := file.Position(file.LineStart(lineNumber)).Offset
	end := file.Position(file.LineStart(Min(file.LineCount(), lineNumber + 1))).Offset
	length := end - start

	// open the file for reading
	filename := file.Name()
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	// seek to the start offset
	_, err = f.Seek(int64(start), 0)
	if err != nil {
		panic(err)
	}
	// then create a buffer of suitable size to hold the line
	line := make([]byte, length)
	// and read exactly the right amount of bytes from the file
	_, err = io.ReadAtLeast(f, line, length)
	if err != nil {
		panic(err)
	}

	// print the position information about the code line, and the line trimmed off any tabs and newlines
	fmt.Printf("%s: %s\n",
		pkg.Fset.File(n.Pos()).Position(n.Pos()).String(),
		strings.Trim(string(line), "\n\t "))
}

/**
 * returns the smaller of two integers. Needed because Go only natively offers a minimum function for floating
 * point values
 */
func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
