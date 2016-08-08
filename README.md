[![Build Status](https://travis-ci.org/ggaaooppeenngg/safemap.svg?branch=master)](https://travis-ci.org/ggaaooppeenngg/safemap)

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

## apis

1. New{{.TypeKey}}2{{.TypeValue}}SafeMap

       New{{.TypeKey}}2{{.TypeValue}}SafeMap returns a new map.

2. {{.TypeKey}}2{{.TypeValue}}SafeMap.Get(k {{.TypeKey}}) \*{{.TypeValue}}

       Get returns a point of {{.TypeValue}}, it returns nil if not found.

3. {{.TypeKey}}2{{.TypeValue}}SafeMap.Set(k {{.TypeKey}}, v {{.TypeValue}})

       Set sets value v to key k in the map.

3. {{.TypeKey}}2{{.TypeValue}}SafeMap.Update(k {{.TypeKey}}, v {{.TypeValue}}) bool

       Update updates value v to key k, returns false if k not found.

4. {{.TypeKey}}2{{.TypeValue}}SafeMap.Delete(k {{.TypeKey}})

       Delete deletes a key in the map.
