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

func (c *Client) Call(args ...string) *Reply {
    req := NewRequest(c.Connect())
    req.wbuf.Write(format(args...))
    req.wbuf.Flush()
    res := req.Read()
    c.Pool.Push(req.conn)
    return res
}

func (c *Client) Connect() net.Conn {
    conn := c.Pool.Pop()

    if conn == nil {
        conn, _ = NewConn(c.Addr, c.Proto)
    }

    return conn
}

func (c *Client) Pipeline() *Pipeline {
    return &Pipeline{c, NewRequest(c.Connect())}
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
