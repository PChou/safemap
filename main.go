// Package safemap is a package used to generate thread-safe map for general purpose.
package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"text/template"
)

var safeMapTemplate = `// Automatically generated file; DO NOT EDIT
package {{.packageName}}

import (
	"sync"
)

type {{ call .builtinType2UCapital .TypeKey }}2{{.TypeValue}}SafeMap struct {
	m    map[{{.TypeKey}}]{{.TypeValue}}
	lock sync.RWMutex
}

func New{{ call .builtinType2UCapital .TypeKey}}2{{.TypeValue}}SafeMap() *{{.TypeKey}}2{{.TypeValue}}SafeMap {
	return &{{.TypeKey}}2{{.TypeValue}}SafeMap{
		m: make(map[{{.TypeKey}}]{{.TypeValue}}),
	}

}

func (s *{{ call .builtinType2UCapital .TypeKey }}2{{.TypeValue}}SafeMap) Get(k {{.TypeKey}}) *{{.TypeValue}} {
	s.lock.RLock()
	v, ok := s.m[k]
	s.lock.RUnlock()
	if !ok {
		return nil
	}
	return &v
}

func (s *{{ call .builtinType2UCapital .TypeKey }}2{{.TypeValue}}SafeMap) Set(k {{.TypeKey}}, v {{.TypeValue}}) {
	s.lock.Lock()
	s.m[k] = v
	s.lock.Unlock()
}

func (s *{{ call .builtinType2UCapital .TypeKey }}2{{.TypeValue}}SafeMap) Update(k {{.TypeKey}}, v {{.TypeValue}}) bool {
	s.lock.Lock()
	_, ok := s.m[k]
	if !ok {
		s.lock.Unlock()
		return false
	}
	s.m[k] = v
	s.lock.Unlock()
	return true
}

func (s *{{ call .builtinType2UCapital .TypeKey }}2{{.TypeValue}}SafeMap) Delete(k {{.TypeKey}}) {
	s.lock.Lock()
	delete(s.m, k)
	s.lock.Unlock()
}`

func fatal(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func builtinType2UCapital(s string) string {
	switch s {
	case "bool":
		return "Bool"
	case "uint8":
		return "Uint8"
	case "uint16":
		return "Uint16"
	case "uint32":
		return "Uint32"
	case "uint64":
		return "Uint64"
	case "int8":
		return "Int8"
	case "int16":
		return "Int16"
	case "int32":
		return "Int32"
	case "int64":
		return "Int64"
	case "float32":
		return "Float32"
	case "float64":
		return "Float64"
	case "complex64":
		return "Complex64"
	case "complex128":
		return "Complex128"
	case "byte":
		return "Byte"
	case "rune":
		return "Rune"
	case "uint":
		return "Uint"
	case "int":
		return "Int"
	case "uintptr":
		return "Uintptr"
	case "string":
		return "String"
	default:
		return s
	}
}

func main() {
	keyType := flag.String("k", "", "key type")
	valueType := flag.String("v", "", "value type")
	flag.Parse()
	if *keyType == "" {
		fatal("key empty")
	}
	if *valueType == "" {
		fatal("value empty")
	}
	tpl, err := template.New("safemap").Parse(safeMapTemplate)
	if err != nil {
		fatal(err)
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, parser.ParseComments)
	if err != nil {
		fatal(err)
	}
	var packageName string
	for name := range pkgs {
		packageName = name
	}
	if packageName == "" {
		fatal("no package found")
	}
	f, err := os.OpenFile(fmt.Sprintf("%s2%s_safemap.go", *keyType, *valueType), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	err = tpl.Execute(f, map[string]interface{}{
		"TypeKey":              *keyType,
		"TypeValue":            *valueType,
		"builtinType2UCapital": builtinType2UCapital,
		"packageName":          packageName,
	})
	if err != nil {
		fatal(err)
	}
}
