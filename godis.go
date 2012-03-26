// package godis implements a db client for Redis. 
package godis

import (
    "bufio"
    "net"
    "strings"
)

type Client struct {
    Addr  string
    Proto string
    Pool  *ConnPool
}

func NewClient(addr string) *Client {
    if addr == "" {
        addr = "tcp:127.0.0.1:6379"
    }

    na := strings.SplitN(addr, ":", 2)
    return &Client{Addr: na[1], Proto: na[0], Pool: NewConnPool()}
}

func (c *Client) Call(args ...string) (*Reply, error) {
    conn, err := c.Connect()
    defer c.Pool.Push(conn)

    if err != nil {
        return nil, err
    }

    req := NewRequest(conn)

    _, err = req.wbuf.Write(format(args...))

    if err != nil {
        return nil, err
    }

    err = req.wbuf.Flush()

    if err != nil {
        return nil, err
    }

    res := req.Read()

    if res.Err != nil {
        return nil, err
    }

    return res, nil
}

func (c *Client) Connect() (conn net.Conn, err error) {
    conn = c.Pool.Pop()

    if conn == nil {
        conn, err = NewConn(c.Addr, c.Proto)

        if err != nil {
            return nil, err
        }
    }

    return conn, nil
}

func (c *Client) Pipeline() (*Pipeline, error) {
    conn, err := c.Connect()

    if err != nil {
        return nil, err
    }

    return &Pipeline{c, NewRequest(conn)}, nil
}

type Pipeline struct {
    *Client
    req *Request
}

func (p *Pipeline) Call(args ...string) {
    p.req.wbuf.Write(format(args...))
}

type Request struct {
    rbuf *bufio.Reader
    wbuf *bufio.Writer
    conn net.Conn
}

func NewRequest(c net.Conn) *Request {
    return &Request{bufio.NewReader(c), bufio.NewWriter(c), c}
}

// reads a reply for a Request
func (r *Request) Read() *Reply {
    if r.wbuf.Buffered() > 0 {
        r.wbuf.Flush()
    }

    res := Parse(r.rbuf)
    return res
}
