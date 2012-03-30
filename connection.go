package godis

import (
    "net"
    "bufio"
)

type Conn struct {
    rbuf *bufio.Reader
    wbuf *bufio.Writer
    conn net.Conn
}

// reads a reply for a Conn
func (r *Conn) Read() *Reply {
    if r.wbuf.Buffered() > 0 {
        r.wbuf.Flush()
    }

    res := Parse(r.rbuf)
    return res
}

// New connection
func newConn(addr, proto string) (*Conn, error) {
    c, err := net.Dial(proto, addr)

    if err != nil {
        return nil, err 
    }

    return &Conn{bufio.NewReader(c), bufio.NewWriter(c), c}, nil
}
