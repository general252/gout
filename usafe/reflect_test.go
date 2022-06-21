package usafe

import (
	"testing"
)

func TestReflectNew(t *testing.T) {

	type Hello struct {
		Name string
	}

	type World struct {
		h   Hello
		Age int
	}

	var (
		a = &World{}
		b = new(World)
		c = World{}
	)
	t.Logf("\n%#v \n%#v \n%#v", a, b, c)

	if x, err := ReflectNew(a); err != nil {
		t.Error(err)
	} else {
		v := x.(World)
		t.Logf("%#v %#v", x, v)
	}

	if x, err := ReflectNew(b); err != nil {
		t.Error(err)
	} else {
		v := x.(World)
		t.Logf("%#v %#v", x, v)
	}

	if x, err := ReflectNew(c); err != nil {
		t.Error(err)
	} else {
		v := x.(World)
		t.Logf("%#v %#v", x, v)
	}
}
