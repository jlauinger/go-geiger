# go-geiger

[![build](https://github.com/jlauinger/go-geiger/workflows/build/badge.svg)](https://github.com/jlauinger/go-geiger/actions/)

![go-geiger logo](https://user-images.githubusercontent.com/1872086/88236443-55c25980-cc7d-11ea-9e81-15c28a8e7daa.png)

Find and count `unsafe` usages in Go packages and their dependencies.


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
      --filter-context string   Count only lines of requested context type (all,variable,parameter,assignment,call,other). Default all (default "all")
      --filter-match string     Count only lines of requested match type (all,pointer,sizeof,offsetof,alignof,sliceheader,stringheader,uintptr). Default pointer (default "pointer")
  -h, --help                    help for geiger
  -q, --hide-stats              Hide statistics table, print only code. --show-code needs to be set manually
      --include-std             Show / include Golang stdlib packages
  -l, --link                    Print link to pkg.go.dev instead of package name
  -d, --max-depth int           Maximum transitive import depth (default 10)
      --show-code               Print the code lines with unsafe usage
      --show-only-once          Do not repeat packages, show them only once and abbreviate further imports (default true)
  -v, --verbose                 Show usage counts by different usage types
```


## Unsafe Match Types

By default, `go-geiger` will count only `unsafe.Pointer` usages. By setting the `--filter-match` argument to one of
`sizeof`, `offsetof`, `alignof`, `sliceheader`, `stringheader`, `uintptr`, or `all`, you can also use `go-geiger` to
find usages of `unsafe.Sizeof`, `unsafe.Offsetof`, `unsafe.Alignof`, `reflect.SliceHeader`, `reflect.StringHeader`,
`uintptr`, or all of them at the same time.


## Unsafe Context Types

Using the `--verbose` argument, you can instruct `go-geiger` to show individual counts for different usage contexts
of unsafe. `go-geiger` distinguishes between the following:

Variable

```go
var x unsafe.Pointer
```

Parameter

```go
func foo(x unsafe.Pointer) {}
```

Assignment

```go
x := unsafe.Pointer(&y)
```

Call

```go
x := unsafe.Pointer(&y)
foo(x)
```

Other, which includes everything that doesn't fall under the first four.

Use the `--filter-context` argument to filter counting to a specific context type. You can use `variable`, `parameter`,
`assignment`, `call`, `other`, or `all`.


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

Run the tests with `go test`.


## License

Licensed under the MIT License (the "License"). You may not use this project except in compliance with the License. You
may obtain a copy of the License [here](https://opensource.org/licenses/MIT).

Copyright 2020 Johannes Lauinger

This tool has been developed as part of my Master's thesis at the
[Software Technology Group](https://www.stg.tu-darmstadt.de/stg/homepage.en.jsp) at TU Darmstadt.

