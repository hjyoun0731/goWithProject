package trace

import (
	"fmt"
	"io"
)

// Tracer 는 코드 전체에서 이벤트를 추적할 수 있는
// 객체를 설명하는 인터페이스다
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// New return Tracer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
