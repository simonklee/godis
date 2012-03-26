package godis

import (
    "net"
)

var MaxConnections = 10

type ConnPool struct {
    free chan net.Conn
}

func NewConnPool() *ConnPool {
    p := ConnPool{make(chan net.Conn, MaxConnections)}

    for i := 0; i < MaxConnections; i++ {
        p.free <- nil
    }

    return &p
}

func (p *ConnPool) Pop() net.Conn {
    return <-p.free
}

func (p *ConnPool) Push(c net.Conn) {
    p.free <- c
}
