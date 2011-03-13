package godis

import (
    "testing"
    "net"
    "log"
)

type C struct {
    Port int
    Name string
    p    *Pool
}

func TestC(t *testing.T) {
    var c C
    log.Println(c.Port)
    if c.p == nil {
        log.Println("nil")
    }

    if c.Name == "" {
        log.Println("empty")
    }

    c.p = NewPool(defaultAddr)
}

func TestPool(t *testing.T) {
    p := NewPool(defaultAddr)

    for i := 0; i < 10; i++ {
        c, err := p.Pop()
        if err != nil {
            t.Errorf("Connection Error", err)
        }

        go func(c *net.TCPConn) {
            p.Push(c)
        }(c)
    }
}
