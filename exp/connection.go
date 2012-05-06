package redis

import (
    "github.com/simonz05/godis/bufin"
    "net"
)

var ConnSum = 0

type Connection interface {
    Write(args ...interface{}) error
    Read() (*Reply, error)
    Close() error
    Sock() net.Conn
}

// Conn implements the Connection interface. 
type Conn struct {
    rbuf *bufin.Reader
    c    net.Conn
}

// NewConn expects a network address and protocol.
// 
//     NewConn("127.0.0.1:6379", "tcp")
// 
// or for a unix domain socket
// 
//     NewConn("/path/to/redis.sock", "unix")
//
// NewConn then returns a Conn struct which implements the Connection
// interface. It's easy to use this interface to create your own
// redis client or to simply talk to the redis database. 
func NewConn(addr, proto string) (*Conn, error) {
    c, err := net.Dial(proto, addr)

    if err != nil {
        return nil, err
    }

    ConnSum++
    return &Conn{bufin.NewReader(c), c}, nil
}

// Read reads one reply of the socket connection. If there is no reply waiting
// this method will block.
// Returns either an error or a pointer to a Reply object.
func (c *Conn) Read() (*Reply, error) {
    reply := Parse(c.rbuf)

    if reply.Err != nil {
        return nil, reply.Err
    }

    return reply, nil
}

// Write accepts any redis command and arbitrary list of arguments.
// 
//     Write("SET", "counter", 1)
//     Write("INCR", "counter")
//
// Write might return a net.Conn.Write error
func (c *Conn) Write(args ...interface{}) error {
    _, e := c.c.Write(format(args...))

    if e != nil {
        return e
    }

    return nil
}

// Close is a simple helper method to close socket connection.
func (c *Conn) Close() error {
    return c.c.Close()
}

// Sock returns the underlying net.Conn. You can use this connection as you
// wish. An example could be to set a r/w deadline on the connection.
//
//      Sock().SetDeadline(t) 
func (c *Conn) Sock() net.Conn {
    return c.c
}
