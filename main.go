// Package safemap is a package used to generate thread-safe map for general purpose.
// Run safemap -k K -v V will generate a K2VSafeMap in K2V_safemap.go for map implementation code,
// to avoid conflicts you can use -n to specify a namespace to generate a namespaceSafeMap in
// namespace_safemap.go. and run go doc you can get document of it.
package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"sort"
	"strings"
	"text/template"
)

var safeMapTemplate = `package {{.PackageName}}

// Automatically generated file; DO NOT EDIT

import (
	"sync"{{ .Imports }}
)

// {{ .Namespace }}SafeMap is a thread-safe map mapping from {{ .TypeKey }} to {{ .TypeValue }}.
type {{ .Namespace }}SafeMap struct {
	m    map[{{.TypeKey}}]{{.TypeValue}}
	lock sync.RWMutex
}

// New{{ .Namespace }}SafeMap returns a new {{ .Namespace }}SafeMap.
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

func formatImports(imports string) string {
	if imports != "" {
		importedPackages := strings.Split(imports, ",")
		sort.Strings(importedPackages)
		for i, s := range importedPackages {
			// wrap package with quote
			importedPackages[i] = "\"" + s + "\""
		}
		imports = strings.Join(importedPackages, "\n\t")
		if len(imports) > 0 {
			imports = "\n\n\t" + imports
		}
	}
	return imports
}

func main() {
	var (
		// flags
		keyType   string
		valueType string
		nameSpace string
		imports   string

		// default package name is main
		packageName = "main"
	)
	flag.StringVar(&keyType, "k", "", "key type")
	flag.StringVar(&valueType, "v", "", "value type")
	flag.StringVar(&nameSpace, "n", "", "namespace")
	flag.StringVar(&imports, "i", "", "imported packages, comma seperated")
	flag.Parse()
	// initiiate paramaters
	if keyType == "" {
		fatal("key empty")
	}
	if valueType == "" {
		fatal("value empty")
	}
	if nameSpace == "" {
		nameSpace = fmt.Sprintf("%s2%s", keyType, valueType)
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
	for name := range pkgs {
		if name != "" {
			packageName = name
			break
		}
	}
	f, err := os.OpenFile(strings.ToLower(fmt.Sprintf("%s_safemap.go", strings.ToLower(nameSpace))), os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fatal(err)
	}
	defer f.Close()
	imports = formatImports(imports)
	err = tpl.Execute(f, map[string]interface{}{
		"TypeKey":     keyType,
		"TypeValue":   valueType,
		"Namespace":   nameSpace,
		"PackageName": packageName,
		"Imports":     imports,
	})
	if err != nil {
		fatal(err)
	}
}
