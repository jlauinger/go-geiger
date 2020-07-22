# go-geiger

[![build](https://github.com/jlauinger/go-geiger/workflows/build/badge.svg)](https://github.com/jlauinger/go-geiger/actions/)

Find and count `unsafe.Pointer` usages in Go packages and their dependencies.

![go-geiger logo](https://user-images.githubusercontent.com/1872086/88235953-3b3bb080-cc7c-11ea-8f99-ca4064ac83f1.png)

*It's dangerous to Go alone. Take \*this!*


## Output example

```
go-geiger -v github.com/jlauinger/go-geiger
```

![go-geiger output example](https://user-images.githubusercontent.com/1872086/88232276-dc733880-cc75-11ea-8081-bab01106390b.png)


## What is the benefit?

A Go package can avoid the restrictions Go normally sets on pointer use by using `unsafe.Pointer`. This can be valuable
to improve efficiency or even necessary e.g. to interact with C code or syscalls.

However, developers must use extreme caution with `unsafe.Pointer` because mistakes can happen quickly and
[they can introduce serious vulnerabilities](https://dev.to/jlauinger/exploitation-exercise-with-unsafe-pointer-in-go-information-leak-part-1-1kga)
such as use-after-free, buffer reuses or buffer overflows.

Since usages of `unsafe.Pointer` can be introduced through dependencies (and dependencies of dependencies), it is necessary
to audit not only the project code but also its dependencies, or at least know who one needs to trust.

`go-geiger` helps developers to quickly identify which packages in the import tree of a Go package use `unsafe.Pointer`, so that developers
can focus auditing efforts onto those, or decide to switch libraries for one that does not use `unsafe.Pointer`.


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
  -f, --filter string    Print only lines of requested type (variable,parameter,assignment,call,other). You need to specify --show-code also. (default "all")
  -h, --help             help for go-geiger
  -q, --hide-stats       Hide statistics table, print only code. --show-code needs to be set manually
      --include-std      Show / include Golang stdlib packages
  -l, --link             Print link to pkg.go.dev instead of package name
  -d, --max-depth int    Maximum transitive import depth (default 10)
      --show-code        Print the code lines with unsafe usage
      --show-only-once   Do not repeat packages, show them only once and abbreviate further imports (default true)
  -v, --verbose          Show usage counts by different usage types

```


## Dependency management

If your project uses Go modules and a `go.mod` file, `go-geiger` will fetch all dependencies automatically before it
analyzes them. It behaves exactly like `go build` would.

If you use a different form of dependency management, e.g. manual `go get`, `go mod vendor` or anything else, you need
to run your dependency management before running `go-geiger` in order to have all dependencies up to date before
analysis.


## Related work

`go-geiger` is inspired by [Cargo Geiger](https://github.com/rust-secure-code/cargo-geiger), a similar tool to find unsafe
code blocks in Rust programs and their dependencies.

[jlauinger/go-unsafepointer-poc](https://github.com/jlauinger/go-unsafepointer-poc) contains proof of concepts for exploiting
vulnerabilities caused by misuse of `unsafe.Pointer`. I also wrote a [blog post series](https://dev.to/jlauinger/exploitation-exercise-with-unsafe-pointer-in-go-information-leak-part-1-1kga)
on the specific problems and vulnerabilities.

[go-safer](https://github.com/jlauinger/go-safer) is a Go linter tool that can help to identify two common and dangerous usage
patterns of `unsafe.Pointer`.


## Development

To get the source code and compile the binary, run this:

```
$ git clone https://github.com/jlauinger/go-geiger
$ cd go-geiger
$ go build
```


## License

Licensed under the MIT License (the "License"). You may not use this project except in compliance with the License. You
may obtain a copy of the License [here](https://opensource.org/licenses/MIT).

Copyright 2020 Johannes Lauinger

This tool has been developed as part of my Master's thesis at the
[Software Technology Group](https://www.stg.tu-darmstadt.de/stg/homepage.en.jsp) at TU Darmstadt.

