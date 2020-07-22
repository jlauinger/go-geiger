package counter

import "golang.org/x/tools/go/packages"

var standardPackages = make(map[string]struct{})

func init() {
	pkgs, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}

	for _, p := range pkgs {
		standardPackages[p.PkgPath] = struct{}{}
	}
}

func isStandardPackage(pkg *packages.Package) bool {
	_, ok := standardPackages[pkg.PkgPath]
	return ok
}

