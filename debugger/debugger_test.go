package debugger

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

var runExample = flag.Bool("example", false, "run debugger on example program")

func TestMain(m *testing.M) {
	flag.Parse()
	if *runExample {
		Example()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestGetPredeclared(t *testing.T) {
	t.Parallel()

	predeclared := starlark.StringDict{"secret": starlark.String("secret")}

	thread := &starlark.Thread{}
	file, err := starlark.ExecFileOptions(&syntax.FileOptions{}, thread, "file", "def f(): return 0", predeclared)
	if err != nil {
		t.Fatal(err)
	}
	f := file["f"].(*starlark.Function)

	predeclared2 := getPredeclared(f)
	if predeclared2["secret"] != predeclared["secret"] {
		t.Errorf("wrong predeclared")
	}
}

type (
	fakeMapValue    map[string]string
	fakeFuncValue   func()
	fakeStructValue struct {
		f func()
		a int
	}
)

func TestSafeEq(t *testing.T) {
	t.Parallel()

	t.Run("Map", func(t *testing.T) {
		t.Parallel()

		a, b, nilMap := fakeMapValue{}, fakeMapValue{}, fakeMapValue(nil)
		assert.False(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(b)))
		assert.True(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(a)))
		assert.False(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(nilMap)))
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		f := func() {}
		a, b, nilFunc := fakeFuncValue(f), fakeFuncValue(f), fakeFuncValue(nil)
		assert.True(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(b)))
		assert.True(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(a)))
		assert.False(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(nilFunc)))
		// Remember safeEq errs on the side of returning true, so don't test anything more complex than this.
	})

	t.Run("Struct", func(t *testing.T) {
		t.Parallel()

		f := func() {}
		a, b, nilFunc := fakeStructValue{f: f, a: 1}, fakeStructValue{f: f, a: 2}, fakeStructValue{f: nil, a: 1}
		assert.False(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(b)))
		assert.True(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(a)))
		assert.False(t, safeEq(reflect.ValueOf(a), reflect.ValueOf(nilFunc)))
	})
}
