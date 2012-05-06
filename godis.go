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
    conn, err := c.Connect()
    defer c.Push(conn)

    if err != nil {
        return nil, err
    }

    conn.Write(args...)

    if err != nil {
        return nil, err
    }

    return conn.Read()
}

// pop a connection from pool 
func (c *Client) Connect() (conn Connection, err error) {
    conn = c.pool.pop()

    if conn == nil {
        conn, err = NewConn(c.Addr, c.Proto)

        if err != nil {
            return nil, err
        }
    }

    return conn, nil
}

// return connection to pool
func (c *Client) Push(conn Connection) {
    c.pool.push(conn)
}

func (c *Client) AsyncClient() *AsyncClient {
    return &AsyncClient{c, bytes.NewBuffer(make([]byte, 0, 1024*16)), nil, 0}
}

type AsyncClient struct {
    *Client
    buf    *bytes.Buffer
    conn   Connection
    queued int
}

func NewAsyncClient(addr string) *AsyncClient {
    return &AsyncClient{
        NewClient(addr),
        bytes.NewBuffer(make([]byte, 0, 1024*16)),
        nil,
        0,
    }
}

func (ac *AsyncClient) Call(args ...interface{}) (err error) {
    _, err = ac.buf.Write(format(args...))
    ac.queued++
    return err
}

func (ac *AsyncClient) Poll() (*Reply, error) {
    if ac.conn == nil {
        conn, e := NewConn(ac.Addr, ac.Proto)

        if e != nil {
            return nil, e
        }

        ac.conn = conn
    }

    if ac.buf.Len() > 0 {
        _, err := ac.buf.WriteTo(ac.conn.Sock())

        if err != nil {
            return nil, err
        }
    }

    reply, e := ac.conn.Read()
    ac.queued--
    return reply, e
}

func (ac *AsyncClient) Close() {
    ac.conn.Close()
    ac.conn = nil
}
