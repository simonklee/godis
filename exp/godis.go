// Package redis implements a db client for Redis.
//
// Connection interface
//
// The Connection interface is a very simple interface to Redis. The Conn
// struct implements this interface and can be used to write commands and read
// replies from Redis.
//
// The Connection interface is used to implement the Client and AsyncClient.
// Unless you like to implment your own client, use either of them instead of a
// single connection.
//
// Client
//
// The Client implements one method; Call(). This writes your command to the
// database, then reads the subsequent reply and returns it to you. 
//
// The Client struct also has a pool of connections so it's safe to use a
// client in a concurrent context. You can create one client for your entire
// program and share it between go routines.
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
// single method Call(), but this method does not return any reply, only an
// error or nil. 
//
//      c := redis.NewAsyncClient("tcp:127.0.0.1:6379")
//      c.Call("SET", "foo", 1)
//      c.Call("GET", "foo")
//
// When we send our command and arguments to the Call() method nothing is sent
// to the Redis server. To get the reply for our commands from Redis we use the
// Read() method. Read sends any buffered commands to the Redis server, and
// then reads one reply. Subsequent calls to Read will return more replies or
// block if there are none.
//
//      // reply from SET 
//      reply, _ := c.Read()
//
//      // reply from GET
//      reply, _ = c.Read()
//
//      println(reply.Elem.Int()) // prints 1
// 
// Due to the nature of how the AsyncClient works, it's not safe to share it
// between go routines.
package redis

import (
    "bytes"
    "errors"
    "strings"
)

// Client implements a Redis client which handles connections to the database
// in a pool. The size of the pool can be adjusted with by setting the
// MaxConnections variable before creating a client.
type Client struct {
    Addr     string
    Proto    string
    Db       int
    Password string
    pool     *connPool
}

// NewClient expects a addr like "tcp:127.0.0.1:6379"
// It returns a new *Client.
func NewClient(addr string, db int, password string) *Client {
    if addr == "" {
        addr = "tcp:127.0.0.1:6379"
    }

    na := strings.SplitN(addr, ":", 2)
    return &Client{na[1], na[0], db, password, newConnPool()}
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
        conn, err = NewConn(c.Addr, c.Proto, c.Db, c.Password)

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
func NewAsyncClient(addr string, db int, password string) *AsyncClient {
    return &AsyncClient{
        NewClient(addr, db, password),
        bytes.NewBuffer(make([]byte, 0, 1024*16)),
        nil,
        0,
    }
}

// Call appends a command to the write buffer or returns an error.
func (ac *AsyncClient) Call(args ...interface{}) {
	// note: bytes.Buffer.Write never returns an error
    _, _ = ac.buf.Write(format(args...))
    ac.queued++
}

// Issue a synchronous call on an async connection.
// This is useful for issuing WATCH commands and doing
// tests before issueing a MULTI command.
// This fails if there are queued commands already.
func (ac *AsyncClient) SyncCall(args ...interface{}) (*Reply, error) {
    if ac.Queued() > 0 {
        return nil, errors.New("Cannot call SyncCall with non-empty queue")
    }

    ac.Call(args...)
    return ac.Read()
}

// Read does three things. 
// 
//      1) Open connection to Redis server, if there is none.
//      2) Write any buffered commands to the server.
//      3) Try to read a reply from the server, or block on read.
//
// Read returns a Reply or error.
func (ac *AsyncClient) Read() (*Reply, error) {
    if ac.conn == nil {
        conn, e := NewConn(ac.Addr, ac.Proto, ac.Db, ac.Password)

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

func (ac *AsyncClient) Queued() int {
    return ac.queued
}

func (ac *AsyncClient) ReadAll() ([]*Reply, error) {
    replies := make([]*Reply, 0, ac.queued)

    for ac.Queued() > 0 {
        r, e := ac.Read()

        if e != nil {
            return nil, e
        }

        replies = append(replies, r)
    }

    return replies, nil
}

// The AsyncClient will only open one connection. This is not automatically
// closed, so to close it we need to call this method.
func (ac *AsyncClient) Close() {
    ac.conn.Close()
    ac.conn = nil
}
