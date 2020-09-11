package counter

import "golang.org/x/tools/go/packages"

// a hash set of all standard packages, which is initialized when this package is imported by the init function
var standardPackages = make(map[string]struct{})

/**
 * initializes the set of standard packages for later querying
 */
func init() {
	// all Go standard library packages can be found by loading the special std path
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	// go through all the packages and add their name to the set
	for _, p := range pkgs {
		standardPackages[p.PkgPath] = struct{}{}
	}
}

/**
 * returns true if the specified package is part of the Go standard library
 */
func isStandardPackage(pkg *packages.Package) bool {
	// check if the package path is included in the prepared set of standard library packages
	_, ok := standardPackages[pkg.PkgPath]
	return ok
}

