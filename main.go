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

type {{.TypeKey}}2{{.TypeValue}}SafeMap struct {
	m    map[{{.TypeKey}}]{{.TypeValue}}
	lock sync.RWMutex
}

func New{{.TypeKey}}2{{.TypeValue}}SafeMap() *{{.TypeKey}}2{{.TypeValue}}SafeMap {
	return &{{.TypeKey}}2{{.TypeValue}}SafeMap{
		m: make(map[{{.TypeKey}}]{{.TypeValue}}),
	}

}

func (s *{{.TypeKey}}2{{.TypeValue}}SafeMap) Get(k {{.TypeKey}}) *{{.TypeValue}} {
	s.lock.RLock()
	v, ok := s.m[k]
	s.lock.RUnlock()
	if !ok {
		return nil
	}
	return &v
}

func (s *{{.TypeKey}}2{{.TypeValue}}SafeMap) Set(k {{.TypeKey}}, v {{.TypeValue}}) {
	s.lock.Lock()
	s.m[k] = v
	s.lock.Unlock()
}

func (s *{{.TypeKey}}2{{.TypeValue}}SafeMap) Update(k {{.TypeKey}}, v {{.TypeValue}}) bool {
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

func (s *{{.TypeKey}}2{{.TypeValue}}SafeMap) Delete(k {{.TypeKey}}) {
	s.lock.Lock()
	delete(s.m, k)
	s.lock.Unlock()
}`

func fatal(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
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
	err = tpl.Execute(f, map[string]string{
		"TypeKey":     *keyType,
		"TypeValue":   *valueType,
		"packageName": packageName,
	})
	if err != nil {
		fatal(err)
	}
}
