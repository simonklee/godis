package godis

import (
    "testing"
    "time"
)

func TestPool(t *testing.T) {
    p := NewConnPool()

    for i := int(MaxConnections) + 1; i >= 0; i-- {
        c := p.Pop()

        go func() {
            time.Sleep(1e+4)
            p.Push(c)
        }()
    }
}
