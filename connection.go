package godis

import (
    "net"
)

var ConnSum = 0

type Conn struct {
    rbuf *Reader
    c net.Conn
}

// reads a reply for a Conn
func (c *Conn) Read() *Reply {
    return Parse(c.rbuf)
}

// New connection
func newConn(addr, proto string) (*Conn, error) {
    c, err := net.Dial(proto, addr)

    if err != nil {
        return nil, err
    }

    ConnSum++
    return &Conn{NewReader(c), c}, nil
}
