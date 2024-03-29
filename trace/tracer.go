package trace

import (
	"fmt"
	"io"
)

//Tracer is the interface that describes an object capable of
//tracing evenets through code.

type Tracer interface {
	Trace(...interface{})
}

func New(w io.Writer) Tracer {
	//	return nil
	return &tracer{out: w}
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

//Off creates a Tracer that will ignore calls to Trace
func Off() Tracer {
	return &nilTracer{}
}
