package godis

import (
    "bufio"
    "net"
)

type Conn struct {
    rbuf *bufio.Reader
    wbuf *bufio.Writer
    conn net.Conn
}

// reads a reply for a Conn
func (c *Conn) Read() *Reply {
    if c.wbuf.Buffered() > 0 {
        c.wbuf.Flush()
    }

    return Parse(r.rbuf)
}

// New connection
func newConn(addr, proto string) (*Conn, error) {
    c, err := net.Dial(proto, addr)

    if err != nil {
        return nil, err
    }

    return &Conn{bufio.NewReader(c), bufio.NewWriter(c), c}, nil
}
