// package godis implements a client for Redis with support
// for all commands and features such as transactions and
// pubsub.
package godis

import (
    "errors"
    "fmt"
    "io"
    "log"

    "strings"
)

type ReaderWriter interface {
    write(b []byte) (*conn, error)
    read(c *conn) *Reply
    sync() *Sync
}

type Client struct {
    Rw ReaderWriter
}

type PipeClient struct {
    *Client
}

type Sync struct {
    Addr     string
    Db       int
    Password string
    net      string
    pool     *pool
}

type Pipe struct {
    *Sync
    conn        *conn
    appendMode  bool
    transaction bool
    replyCount  int
}

type Sub struct {
    c          *Sync
    conn       *conn
    subscribed bool
    Messages   chan *Message
}

func New(netaddr string, db int, password string) *Client {
    return &Client{newSync(netaddr, db, password)}
}

// New returns a new Sync given a net address, redis db and password.
// nettaddr should be formatted using "net:addr", where ":" is acting as a
// separator. E.g. "unix:/path/to/redis.sock", "tcp:127.0.0.1:12345" or use an
// empty string for redis default.
func newSync(netaddr string, db int, password string) *Sync {
    if netaddr == "" {
        netaddr = "tcp:127.0.0.1:6379"
    }

    na := strings.SplitN(netaddr, ":", 2)

    return &Sync{Addr: na[1], Db: db, Password: password, net: na[0], pool: newPool()}
}

// PipeClient include support for MULTI/EXEC operations. 
// It implements Exec() which executes all buffered
// commands. Set transaction to true to wrap buffered commands inside
// MULTI .. EXEC block.
func NewPipeClient(netaddr string, db int, password string, transaction bool) *PipeClient {
    s := newSync(netaddr, db, password)
    p := &Pipe{s, nil, true, transaction, 0}
    c := &Client{p}
    return &PipeClient{c}
}

// Uses the connection settings from a existing client to create a new PipeClient
func NewPipeClientFromClient(c *Client, transaction bool) *PipeClient {
    s := c.Rw.sync()
    netaddr := s.net + ":" + s.Addr
    return NewPipeClient(netaddr, s.Db, s.Password, transaction)
}

func (p *PipeClient) pipe() *Pipe {
    v, _ := p.Rw.(*Pipe)
    return v
}

func NewSub(addr string, db int, password string) *Sub {
    return &Sub{c: newSync(addr, db, password)}
}

// rw interface 

func (c *Sync) read(conn *conn) *Reply {
    r := conn.readReply()

    if r.Err == io.EOF {
        conn = nil
    }

    c.pool.push(conn)
    return r
}

func (c *Sync) write(cmd []byte) (conn *conn, err error) {
    if conn, err = c.getConn(); err != nil {
        return nil, err
    }

    if _, err = conn.w.Write(cmd); err != nil {
        c.pool.push(conn)
        return nil, err
    }

    conn.w.Flush()
    return conn, err
}

func (c *Sync) sync() *Sync {
    return c
}

// extra methods on sync 
func (c *Sync) getConn() (*conn, error) {
    cc := c.pool.pop()

    if cc != nil {
        return cc, nil
    }

    return newConn(c.net, c.Addr, c.Db, c.Password)
}

// pipe interface implementation
func (p *Pipe) read(conn *conn) *Reply {
    if p.appendMode {
        return &Reply{}
    }

    if p.conn.w.Buffered() > 0 {
        if logCmd {
            log.Printf("%d bytes were written to socket\n", p.conn.w.Buffered())
        }
        p.conn.w.Flush()
    }

    reply := conn.readReply()

    if p.count() == 0 {
        p.free()
    }

    return reply
}

func (p *Pipe) write(cmd []byte) (*conn, error) {
    var err error

    if p.conn == nil {
        if c, err := p.getConn(); err != nil {
            return nil, err
        } else {
            p.conn = c
        }
    }

    if p.transaction && p.replyCount == 0 {
        p.replyCount++
        p.conn.w.Write(buildCmd([][]byte{[]byte("MULTI")}))
    }

    if _, err = p.conn.w.Write(cmd); err != nil {
        p.free()
        return nil, err
    }

    p.appendMode = true
    p.replyCount++
    return p.conn, nil
}

// read a reply from the socket if we are expecting it.
func (p *Pipe) getReply() *Reply {
    if p.count() == 0 {
        p.appendMode = true
        p.transaction = false
        return &Reply{Err: errors.New("No replies expected from conn")}
    }

    p.replyCount--
    p.appendMode = false
    return p.read(p.conn)
}

// retrieve the number of replies available
func (p *Pipe) count() int {
    return p.replyCount
}

func (p *Pipe) free() {
    p.pool.push(p.conn)
    p.conn = nil
    p.appendMode = true
}

func (s *Sub) read(conn *conn) *Reply {
    return s.conn.readReply()
}

func (s *Sub) write(cmd []byte) (*conn, error) {
    var err error

    if s.conn == nil {
        if c, err := s.c.getConn(); err != nil {
            return nil, err
        } else {
            s.conn = c
        }
    }

    if _, err = s.conn.w.Write(cmd); err != nil {
        s.Close()
        return nil, err
    }

    s.conn.w.Flush()
    return s.conn, nil
}

func (s *Sub) sync() *Sync {
    return s.c
}

func (s *Sub) listen() {
    if s.conn == nil {
        return
    }

    for {
        r := s.read(s.conn)

        if r.Err != nil {
            go s.free()
            return
        }

        if m := r.Message(); m != nil {
            s.Messages <- m
        }
    }
}

func (s *Sub) subscribe() {
    s.subscribed = true
    s.Messages = make(chan *Message, 64)
    go s.listen()
}

// Free the connection and close the chan
func (s *Sub) Close() {
    s.conn.rwc.Close()
}

func (s *Sub) free() {
    s.conn = nil
    s.c.pool.push(nil)
    s.subscribed = false

    close(s.Messages)
}

// Methods which take ReaderWriter interface
func sendGen(rw ReaderWriter, readResp bool, retry int, args [][]byte) (r *Reply) {
    c, err := rw.write(buildCmd(args))
    r = &Reply{conn: c, Err: err}

    defer func() {
        // if connection was closed by the remote host we try to re-run the cmd
        if retry > 0 && r.Err == io.EOF {
            retry--
            r = sendGen(rw, readResp, retry, args)
        }
    }()

    if r.Err != nil {
        return
    }

    if readResp {
        return rw.read(c)
    }

    return
}

// writes a command a and returns single the Reply object
func Send(rw ReaderWriter, args ...[]byte) *Reply {
    return sendGen(rw, true, MaxClientConn, args)
}

// uses reflection to create a bytestring of the name and args parameters, 
// then calls Send()
func SendIface(rw ReaderWriter, name string, args ...interface{}) *Reply {
    buf := make([][]byte, len(args)+1)
    buf[0] = []byte(name)

    for i, arg := range args {
        switch v := arg.(type) {
        case []byte:
            buf[i+1] = v
        case string:
            buf[i+1] = []byte(v)
        default:
            buf[i+1] = []byte(fmt.Sprint(arg))
        }
    }

    return sendGen(rw, true, MaxClientConn, buf)
}

func strToBytes(name string, args []string) [][]byte {
    buf := make([][]byte, len(args)+1)
    buf[0] = []byte(name)

    for i, arg := range args {
        buf[i+1] = []byte(arg)
    }
    return buf
}

func appendSendStr(rw ReaderWriter, name string, args ...string) *Reply {
    buf := strToBytes(name, args)
    return sendGen(rw, false, MaxClientConn, buf)
}

// creates a bytestring of the name and args parameters, then calls Send()
func SendStr(rw ReaderWriter, name string, args ...string) *Reply {
    buf := strToBytes(name, args)
    return sendGen(rw, true, MaxClientConn, buf)
}
