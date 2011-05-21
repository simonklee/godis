// godis package implements a client for Redis. It supports all redis commands
// and common features such as pipelines.
package godis

import (
    "fmt"
    "log"
    "os"
    "strings"
)

type ReaderWriter interface {
    write(b []byte) (*conn, os.Error)
    read(c *conn) *Reply
}

type Client struct {
    Addr     string
    Db       int
    Password string
    net      string
    pool     *pool
}

func sendGen(rw ReaderWriter, readResp bool, retry int, args [][]byte) (r *Reply) {
    c, err := rw.write(buildCmd(args))
    r = &Reply{conn: c, Err: err}

    defer func() {
        // if connection was closed by the remote host we try to re-run the cmd
        if retry > 0 && r.Err == os.EOF {
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

// New returns a new Client given a net address, redis db and password.
// nettaddr should be formatted using "net:addr", where ":" is acting as a
// separator. E.g. "unix:/path/to/redis.sock", "tcp:127.0.0.1:12345" or use an
// empty string for redis default.
func New(netaddr string, db int, password string) *Client {
    if netaddr == "" {
        netaddr = "tcp:127.0.0.1:6379"
    }

    na := strings.Split(netaddr, ":", 2)

    return &Client{Addr: na[1], Db: db, Password: password, net: na[0], pool: newPool()}
}

func (c *Client) getConn() (*conn, os.Error) {
    cc := c.pool.pop()

    if cc != nil {
        return cc, nil
    }

    return newConn(c.net, c.Addr, c.Db, c.Password)
}

func (c *Client) read(conn *conn) *Reply {
    r := conn.readReply()

    if r.Err == os.EOF {
        conn = nil
    }

    c.pool.push(conn)
    return r
}

func (c *Client) write(cmd []byte) (conn *conn, err os.Error) {
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

type Pipe struct {
    *Client
    conn       *conn
    appendMode bool
    replyCount int
}

// Pipe implements the ReaderWriter interface, can be used with all commands.
// Currently its not possible to use a Pipe object in a concurrent context.
func NewPipe(addr string, db int, password string) *Pipe {
    return &Pipe{New(addr, db, password), nil, true, 0}
}

func NewPipeFromClient(c *Client) *Pipe {
    return &Pipe{c, nil, true, 0}
}

// read a reply from the socket if we are expecting it.
func (p *Pipe) GetReply() *Reply {
    if p.Count() > 0 {
        p.replyCount--
        p.appendMode = false
    } else {
        p.appendMode = true
        return &Reply{Err: os.NewError("No replies expected from conn")}
    }

    return p.read(p.conn)
}

// retrieve the number of replies available
func (p *Pipe) Count() int {
    return p.replyCount
}

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

    if reply.Err != nil || p.Count() == 0 {
        p.free()
    }

    return reply
}

func (p *Pipe) write(cmd []byte) (*conn, os.Error) {
    var err os.Error

    if p.conn == nil {
        if c, err := p.getConn(); err != nil {
            return nil, err
        } else {
            p.conn = c
        }
    }

    if _, err = p.conn.w.Write(cmd); err != nil {
        p.free()
        return nil, err
    }

    p.appendMode = true
    p.replyCount++
    return p.conn, nil
}

func (p *Pipe) free() {
    p.pool.push(p.conn)
    p.conn = nil
    p.appendMode = true
}

type Sub struct {
    c          *Client
    conn       *conn
    subscribed bool
    Messages   chan *Message
}

func NewSub(addr string, db int, password string) *Sub {
    return &Sub{c: New(addr, db, password)}
}

func (s *Sub) read(conn *conn) *Reply {
    return s.conn.readReply()
}

func (s *Sub) write(cmd []byte) (*conn, os.Error) {
    var err os.Error

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
