package debugger

import (
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func Example() {
	// Run this example with
	//
	//	go test -c github.com/stripe/skycfg/debugger && ./debugger.test -example

	prog := `
def g():
    closure = [42]
    def f():
        local_var = 1
        breakpoint()
        print(closure)
    return f

g()()
`

	dbg := New()
	var thread starlark.Thread
	starlark.ExecFileOptions(&syntax.FileOptions{}, &thread, "<prog>", prog, starlark.StringDict{
		"predeclared": starlark.String("value"),
		"breakpoint":  dbg.Breakpoint(),
	})
}
