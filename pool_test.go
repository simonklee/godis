package godis

import (
    "testing"
    "time"
)

func TestPool(t *testing.T) {
    p := newConnPool()

    for i := int(MaxConnections) + 1; i >= 0; i-- {
        c := p.pop()

        go func() {
            time.Sleep(1e+4)
            p.push(c)
        }()
    }
}
