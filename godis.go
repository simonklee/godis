// package godis implements a db client for Redis. 
package godis

import (
    "bytes"
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

func (c *Client) Call(args ...interface{}) (*Reply, error) {
    conn, err := c.connect()
    defer c.pool.push(conn)

    if err != nil {
        return nil, err
    }

    _, err = conn.c.Write(format(args...))

    if err != nil {
        return nil, err
    }

    res := conn.Read()

    if res.Err != nil {
        return nil, res.Err
    }

    return res, nil
}

func (c *Client) connect() (conn *Conn, err error) {
    conn = c.pool.pop()

    if conn == nil {
        conn, err = newConn(c.Addr, c.Proto)

        if err != nil {
            return nil, err
        }
    }

    return conn, nil
}

func (c *Client) Pipeline() *Pipeline {
    return &Pipeline{c, bytes.NewBuffer(make([]byte, 0, 1024*16)), nil}
}

type Pipeline struct {
    *Client
    buf  *bytes.Buffer
    conn *Conn
}

func (p *Pipeline) Call(args ...interface{}) (err error) {
    _, err = p.buf.Write(format(args...))
    return err
}

func (p *Pipeline) Read() (*Reply, error) {
    if p.conn == nil {
        conn, err := p.connect()

        if err != nil {
            return nil, err
        }

        p.conn = conn
    }

    if p.buf.Len() > 0 {
        _, err := p.buf.WriteTo(p.conn.c)

        if err != nil {
            return nil, err
        }
    }

    res := p.conn.Read()

    if res.Err != nil {
        return nil, res.Err
    }

    return res, nil
}
