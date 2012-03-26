package godis

import (
    "net"
)

var MaxConnections = 10

type connPool struct {
    free chan net.Conn
}

func newConnPool() *connPool {
    p := connPool{make(chan net.Conn, MaxConnections)}

    for i := 0; i < MaxConnections; i++ {
        p.free <- nil
    }

    return &p
}

func (p *connPool) pop() net.Conn {
    return <-p.free
}

func (p *connPool) push(c net.Conn) {
    p.free <- c
}
