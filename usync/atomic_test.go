package usync

import (
	"testing"
)

func TestAtomicString(t *testing.T) {
	var s AtomicString
	if s.Get() != "" {
		t.Errorf("want empty, got %s", s.Get())
	}
	s.Set("a")
	if s.Get() != "a" {
		t.Errorf("want a, got %s", s.Get())
	}
	if s.CompareAndSwap("b", "c") {
		t.Errorf("want false, got true")
	}
	if s.Get() != "a" {
		t.Errorf("want a, got %s", s.Get())
	}
	if !s.CompareAndSwap("a", "c") {
		t.Errorf("want true, got false")
	}
	if s.Get() != "c" {
		t.Errorf("want c, got %s", s.Get())
	}
}

func TestAtomicBool(t *testing.T) {
	var b AtomicBool
	if b.Get() != false {
		t.Fatal("must false")
	}

	b.Set(true)

	if b.Get() != true {
		t.Fatal("must true")
	}

	b.Set(false)

	if b.Get() != false {
		t.Fatal("must false")
	}
}
