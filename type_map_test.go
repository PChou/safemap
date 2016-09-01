package main

import (
	"testing"
)

func TestMap(t *testing.T) {
	m := NewNamespaceSafeMap(nil)
	v, ok := m.Get(TypeKey("k"))
	if ok {
		t.Fatalf("get %s\n", v)
	}
	m.Set(TypeKey("k"), TypeValue("v"))
	v, ok = m.Get(TypeKey("k"))
	if ok && v != TypeValue("v") {
		t.Fatalf("get %s\n", v)
	}
	dupMap := m.Dup()
	v, ok = dupMap.Get(TypeKey("k"))
	if ok && v != TypeValue("v") {
		t.Fatalf("get %s\n", v)
	}
	m.Update(TypeKey("k"), TypeValue("v2"))
	v, ok = m.Get(TypeKey("k"))
	if ok && v != TypeValue("v2") {
		t.Fatalf("get %s\n", v)
	}
	ok = m.Update(TypeKey("k1"), TypeValue("v2"))
	if ok {
		t.Fatalf("get %b\n", ok)
	}
	m.Delete(TypeKey("k"))
	v, ok = m.Get(TypeKey("k"))
	if ok && v != "" {
		t.Fatalf("get %s\n", v)
	}
}
