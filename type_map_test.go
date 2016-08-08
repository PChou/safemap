package main

import (
	"testing"
)

func TestMap(t *testing.T) {
	m := NewTypeKey2TypeValueSafeMap()
	v := m.Get(TypeKey("k"))
	if v != nil {
		t.Fatalf("get %s\n", *v)
	}
	m.Set(TypeKey("k"), TypeValue("v"))
	v = m.Get(TypeKey("k"))
	if *v != TypeValue("v") {
		t.Fatalf("get %s\n", *v)
	}
	dupMap := m.Dup()
	v = dupMap.Get(TypeKey("k"))
	if *v != TypeValue("v") {
		t.Fatalf("get %s\n", *v)
	}
	m.Update(TypeKey("k"), TypeValue("v2"))
	v = m.Get(TypeKey("k"))
	if *v != TypeValue("v2") {
		t.Fatalf("get %s\n", *v)
	}
	ok := m.Update(TypeKey("k1"), TypeValue("v2"))
	if ok {
		t.Fatalf("get %b\n", ok)
	}
	m.Delete(TypeKey("k"))
	v = m.Get(TypeKey("k"))
	if v != nil {
		t.Fatalf("get %s\n", *v)
	}
}
