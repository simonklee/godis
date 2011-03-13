package godis

import (
    "net"
    "os"
)

const (
    MaxClientConn = 5
)

type Pool struct {
    pool chan *net.TCPConn
    addr string
}

func NewPool(addr string) *Pool {
    p := Pool{make(chan *net.TCPConn, MaxClientConn), addr}
    for i := 0; i < MaxClientConn; i++ {
        p.pool <- nil
    }
    return &p
}

func (p *Pool) Pop() (c *net.TCPConn, err os.Error) {
    c = <-p.pool
    if c == nil {
        return connect(p.addr)
    }
    return c, nil
}

func (p *Pool) Push(c *net.TCPConn) {
    p.pool <- c
}

// TODO: flush the pool
func (p *Pool) Flush() {

}

func connect(addr string) (c *net.TCPConn, err os.Error) {
    a, err := net.ResolveTCPAddr(addr)
    if err != nil {
        return nil, os.NewError("ResolveAddr error for " + addr)
    }

    c, err = net.DialTCP("tcp", nil, a)
    if err != nil {
        err = os.NewError("Connection error " + a.String())
    }
    return
}
