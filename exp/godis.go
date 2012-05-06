// Package redis implements a db client for Redis.
//
// Connection interface
//
// The Connection interface is a very simple interface to Redis, with a Read
// and Write method. The Conn struct implements the Connection interface and 
// can be used to read and write commands and replies from Redis.
//
// But you don't need to use the Connection interface, because we have something
// simpler. 
//
// Client
//
// The Client implements one method called `Call`. This method first writes
// your command to Redis, then reads the subsequent reply and returns it to
// you. 
//
// The Client struct also has a pool of connections so it's safe to use a
// client in a concurrent context.  You can create on client object for your
// entire program and share it between go routines.
//
//      c := redis.NewClient("tcp:127.0.0.1:6379")
//      reply, e := c.Call("GET", "foo")
//
//      if e != nil {
//          // handle error
//      }
//
//      println(reply.Elem.String())
//
// AsyncClient
// 
// The AsyncClient works exactly like the regular Client, and implements a
// single method `Call`, but this method does not return any reply, only an
// error or nil. 
//
//      c := redis.NewAsyncClient("tcp:127.0.0.1:6379")
//      c.Call("SET", "foo", 1)
//      c.Call("GET", "foo")
//
// When we send our command and arguments to the Call method nothing is sent to
// the Redis server. To get the reply for our commands from Redis we use the
// `Poll` method. Poll sends any buffered commands to the Redis server, and
// then reads one reply. Subsequent calls to Poll will return more replies or
// block if there are none.
//
//      // reply from SET 
//      reply, _ := c.Poll()
//
//      // reply from GET
//      reply, _ = c.Poll()
//
//      println(reply.Elem.Int()) // prints 1
package redis

import (
    "bytes"
    "strings"
)

// Client implements a redis client which handles connections to the database
// in a pool. The size of the pool can be adjusted with by setting the
// MaxConnections variable before calling NewClient.
type Client struct {
    Addr  string
    Proto string
    pool  *connPool
}

// NewClient expects a addr like "tcp:127.0.0.1:6379"
// It returns a new *Client.
func NewClient(addr string) *Client {
    if addr == "" {
        addr = "tcp:127.0.0.1:6379"
    }

    na := strings.SplitN(addr, ":", 2)
    return &Client{Addr: na[1], Proto: na[0], pool: newConnPool()}
}

// Call is the canonical way of talking to Redis. It accepts any 
// Redis command and a arbitrary number of arguments.
// Call returns a Reply object or an error.
func (c *Client) Call(args ...interface{}) (*Reply, error) {
    conn, err := c.connect()
    defer c.pool.push(conn)

    if err != nil {
        return nil, err
    }

    conn.Write(args...)

    if err != nil {
        return nil, err
    }

    return conn.Read()
}

// Pop a connection from pool 
func (c *Client) connect() (conn Connection, err error) {
    conn = c.pool.pop()

    if conn == nil {
        conn, err = NewConn(c.Addr, c.Proto)

        if err != nil {
            return nil, err
        }
    }

    return conn, nil
}

// Use the connection settings from Client to create a new AsyncClient
func (c *Client) AsyncClient() *AsyncClient {
    return &AsyncClient{c, bytes.NewBuffer(make([]byte, 0, 1024*16)), nil, 0}
}

// Async client implements an asynchronous client. It is very similar to Client
// except that it maintains a buffer of commands which first are sent to Redis
// once we explicitly request a reply.
type AsyncClient struct {
    *Client
    buf    *bytes.Buffer
    conn   Connection
    queued int
}

// NewAsyncClient expects a addr like "tcp:127.0.0.1:6379"
// It returns a new *Client.
func NewAsyncClient(addr string) *AsyncClient {
    return &AsyncClient{
        NewClient(addr),
        bytes.NewBuffer(make([]byte, 0, 1024*16)),
        nil,
        0,
    }
}

// Call appends a command to the write buffer or returns an error.
func (ac *AsyncClient) Call(args ...interface{}) (err error) {
    _, err = ac.buf.Write(format(args...))
    ac.queued++
    return err
}

// Poll does three things. 
// 
//      1) Open connection to Redis server, if there is none.
//      2) Write any buffered commands to the server.
//      3) Try to read a reply from the server, or block on read.
//
// Poll returns a Reply or error.
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

// The AsyncClient will only open one connection. This is not automatically
// closed, so to close it we need to call this method.
func (ac *AsyncClient) Close() {
    ac.conn.Close()
    ac.conn = nil
}
