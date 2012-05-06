package redis

import (
    "net"
    "github.com/simonz05/godis/bufin"
)

var ConnSum = 0

type Connection interface {
    Write(args ...interface{}) error
    Read() (*Reply, error)
    Close() error
    Sock() net.Conn
}

type Conn struct {
    rbuf *bufin.Reader
    c    net.Conn
}

// New connection
func NewConn(addr, proto string) (*Conn, error) {
    c, err := net.Dial(proto, addr)

    if err != nil {
        return nil, err
    }

    ConnSum++
    return &Conn{bufin.NewReader(c), c}, nil
}

// read and parse a reply from socket
func (c *Conn) Read() (*Reply, error) {
    reply := Parse(c.rbuf)

    if reply.Err != nil {
        return nil, reply.Err
    }

    return reply, nil
}

// write args to socket
func (c *Conn) Write(args ...interface{}) error {
    _, e := c.c.Write(format(args...))

    if e != nil {
        return e
    }

    return nil
}

// close socket connection
func (c *Conn) Close() error {
    return c.c.Close()
}

// returns the net.Conn for the struct
func (c *Conn) Sock() net.Conn {
    return c.c
}
