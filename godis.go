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
    pool  *connPool
}

func NewClient(addr string) *Client {
    if addr == "" {
        addr = "tcp:127.0.0.1:6379"
    }

    na := strings.SplitN(addr, ":", 2)
    return &Client{Addr: na[1], Proto: na[0], pool: newConnPool()}
}

func (c *Client) Call(args ...string) (*Reply, error) {
    conn, err := c.connect()
    defer c.pool.push(conn)

    if err != nil {
        return nil, err
    }

    req := newRequest(conn)

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
        return nil, res.Err
    }

    return res, nil
}

func (c *Client) connect() (conn net.Conn, err error) {
    conn = c.pool.pop()

    if conn == nil {
        conn, err = newConn(c.Addr, c.Proto)

        if err != nil {
            return nil, err
        }
    }

    return conn, nil
}

func (c *Client) Pipeline() (*Pipeline, error) {
    //TODO: connect at a later stage
    conn, err := c.connect()

    if err != nil {
        return nil, err
    }

    return &Pipeline{c}, nil
}

type Pipeline struct {
    *Client
    req *Request
}

func (p *Pipeline) Call(args ...string) (error) {
    _, err := p.req.wbuf.Write(format(args...))
    return err
}

func (p *Pipeline) Read() (*Reply, error) {
    res := p.req.Read()

    if res.Err != nil {
        return nil, res.Err
    }

    return res, nil
}

type Request struct {
    rbuf *bufio.Reader
    wbuf *bufio.Writer
    conn net.Conn
}

func newRequest(c net.Conn) *Request {
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
