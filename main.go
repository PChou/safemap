// Package safemap is a package used to generate thread-safe map for general purpose.
// Run safemap -k -v will generate a K2V_safemap.go for map implementation code,
// and run go doc you can get document of it.
package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

var safeMapTemplate = `package {{.packageName}}

// Automatically generated file; DO NOT EDIT

import (
	"sync"
)

// {{ .Namespace }}SafeMap is a thread-safe map mapping from
// {{ .TypeKey }} to {{ .TypeValue }}.
type {{ .Namespace }}SafeMap struct {
	m    map[{{.TypeKey}}]{{.TypeValue}}
	lock sync.RWMutex
}

// New{{ .Namespace }}SafeMap returns a new
// {{ .Namespace }}SafeMap.
func New{{ .Namespace }}SafeMap(m map[{{.TypeKey}}]{{.TypeValue}}) *{{ .Namespace }}SafeMap {
	if m == nil {
		m = make(map[{{.TypeKey}}]{{.TypeValue}})
	}
	return &{{ .Namespace }}SafeMap{
		m: m,
	}

}

// Get returns a point of {{.TypeValue}}, it returns nil if not found.
func (s *{{ .Namespace }}SafeMap) Get(k {{.TypeKey}}) ({{.TypeValue}}, bool) {
	s.lock.RLock()
	v, ok := s.m[k]
	s.lock.RUnlock()
	return v, ok
}

// Set sets value v to key k in the map.
func (s *{{ .Namespace }}SafeMap) Set(k {{.TypeKey}}, v {{.TypeValue}}) {
	s.lock.Lock()
	s.m[k] = v
	s.lock.Unlock()
}

// Update updates value v to key k, returns false if k not found.
func (s *{{ .Namespace }}SafeMap) Update(k {{.TypeKey}}, v {{.TypeValue}}) bool {
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

// Delete deletes a key in the map.
func (s *{{ .Namespace }}SafeMap) Delete(k {{.TypeKey}}) {
	s.lock.Lock()
	delete(s.m, k)
	s.lock.Unlock()
}

// Dup duplicates the map to a new struct.
func (s *{{ .Namespace }}SafeMap) Dup() *{{ .Namespace }}SafeMap {
	newMap := New{{ .Namespace }}SafeMap(nil)
	s.lock.Lock()
	for k, v := range s.m {
		newMap.m[k] = v
	}
	s.lock.Unlock()
	return newMap
}`

func fatal(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func main() {
	keyType := flag.String("k", "", "key type")
	valueType := flag.String("v", "", "value type")
	nameSpace := flag.String("n", "", "namespace")
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
	packageName := "main"
	for name := range pkgs {
		packageName = name
	}
	if *nameSpace == "" {
		*nameSpace = fmt.Sprintf("%s2%s", strings.Title(*keyType), *valueType) // TODO: only title builtin types
	}
	f, err := os.OpenFile(strings.ToLower(fmt.Sprintf("%s_safemap.go", *nameSpace)), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fatal(err)
	}
	defer f.Close()

	err = tpl.Execute(f, map[string]interface{}{
		"TypeKey":     *keyType,
		"TypeValue":   *valueType,
		"Namespace":   *nameSpace,
		"packageName": packageName,
	})
	if err != nil {
		fatal(err)
	}
}
