package main

import (
	"testing"
)

func TestImportsFormat(t *testing.T) {
	if formatImports("") != "" {
		t.Fatalf("Get `%s`", formatImports(""))
	}
	if formatImports("a/b/c,d/e/f") != `

	"a/b/c"
	"d/e/f"` {
		t.Fatalf("Get `%s`", formatImports("a/b/c,d/e/f"))
	}
}
