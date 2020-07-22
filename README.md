# go-geiger

Find and count unsafe usages in Go packages and their dependencies.


## Output example

```
go-geiger -v github.com/jlauinger/go-geiger
```

![go-geiger output example](https://user-images.githubusercontent.com/1872086/88232276-dc733880-cc75-11ea-8081-bab01106390b.png)


## Install

To install `go-geiger`, use the following command:

```
go get github.com/stg-tud/thesis-2020-lauinger-code/go-geiger
```

This will install `go-geiger` to `$GOPATH/bin`, so make sure that it is included in your `$PATH` environment variable.


## Usage

Run go-geiger on a package like this:

```
$ go-geiger example/cmd
```

Or supply multiple packages, separated by spaces:

```
$ go-geiger example/cmd example/util strings
```

To check the package in the current directory you can call `go-geiger` without parameters:

```
$ go-geiger
```

Supplying the `--help` flag prints the usage information for `go-geiger`:

```
$ go-geiger --help
```

There are the following flags available:

```
      --dnr         Do not repeat packages (default true)
  -h, --help        help for geiger
      --level int   Maximum indent level (default 10)
      --link        Print link to pkg.go.dev instead of package name
      --show-code   Print the code lines with unsafe usage
      --show-std    Show Goland stdlib packages
```


## Dependency management

If your project uses Go modules and a `go.mod` file, `go-geiger` will fetch all dependencies automatically before it
analyzes them. It behaves exactly like `go build` would.

If you use a different form of dependency management, e.g. manual `go get`, `go mod vendor` or anything else, you need
to run your dependency management before running `go-geiger` in order to have all dependencies up to date before 
analysis.


## Development

To get the source code and compile the binary, run this:

```
$ git clone https://github.com/stg-tud/thesis-2020-lauinger-code
$ cd thesis-2020-lauinger-code/go-geiger
$ go build
```


## License

Licensed under the MIT License (the "License"). You may not use this project except in compliance with the License. You 
may obtain a copy of the License [here](https://opensource.org/licenses/MIT).

Copyright 2020 Johannes Lauinger

This tool has been developed as part of my Master's thesis at the 
[Software Technology Group](https://www.stg.tu-darmstadt.de/stg/homepage.en.jsp) at TU Darmstadt.
