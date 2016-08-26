[![Build Status](https://travis-ci.org/ggaaooppeenngg/safemap.svg?branch=master)](https://travis-ci.org/ggaaooppeenngg/safemap)
[![Go Report Card](https://goreportcard.com/badge/github.com/ggaaooppeenngg/safemap)](https://goreportcard.com/report/github.com/ggaaooppeenngg/safemap)
[![codecov](https://codecov.io/gh/ggaaooppeenngg/safemap/branch/master/graph/badge.svg)](https://codecov.io/gh/ggaaooppeenngg/safemap)


## safemap
An auto-generated thread-safe map package written in golang.

## install

`go get github.com/ggaaooppeenngg/safemap`

## usage

Run `safemap -k Key_type -v Val_type` to generate a file named `Key_type2Val_type.go` in currenty package directory. It will search currenty directory for the definitions of Key\_type and Val\_type, and define a struct `Key_type2Val_typeSafeMap`  in the file.

You can also use [go generate](https://blog.golang.org/generate) to automatically generate the code, put a comment in your code like below and run `go generate`, a generated file will be found.

```
//go:generate safemap -k TypeKey -v TypeValue
type TypeKey string
type TypeValue string

```
