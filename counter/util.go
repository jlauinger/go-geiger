package counter

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"io"
	"os"
	"strings"
)

func getPrintedPackageName(pkg *packages.Package, config Config) string {
	if config.PrintLinkToPkgGoDev && pathLooksLikeAUrl(pkg.PkgPath) {
		return fmt.Sprintf("https://pkg.go.dev/%s", pkg.PkgPath)
	} else {
		return pkg.PkgPath
	}
}

func pathLooksLikeAUrl(path string) bool {
	components := strings.Split(path, "/")

	if len(components) <= 1 {
		return false
	}

	domain := components[0]

	return strings.Contains(domain, ".")
}

func getImportsCount(pkgs map[string]*packages.Package, config Config) (childCount, stdLibCount int) {
	for _, pkg := range pkgs {
		if isStandardPackage(pkg) {
			stdLibCount++
			if config.ShowStandardPackages {
				childCount++
			}
		} else {
			childCount++
		}
	}
	return
}

func printLine(pkg *packages.Package, n ast.Node) {
	file := pkg.Fset.File(n.Pos())
	lineNumber := file.Position(n.Pos()).Line  // 1-based

	start := file.Position(file.LineStart(lineNumber)).Offset
	end := file.Position(file.LineStart(Min(file.LineCount(), lineNumber + 1))).Offset
	length := end - start

	filename := file.Name()

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	_, err = f.Seek(int64(start), 0)
	if err != nil {
		panic(err)
	}
	line := make([]byte, length)
	_, err = io.ReadAtLeast(f, line, length)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %s\n",
		pkg.Fset.File(n.Pos()).Position(n.Pos()).String(),
		strings.Trim(string(line), "\n\t "))
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
