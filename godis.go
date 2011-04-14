// godis package implements a client for Redis. It supports all redis commands
// and common features such as pipelines.
package godis

import (
    "bytes"
    "fmt"
    "log"
    "net"
    "os"
    "strconv"
)

type ReaderWriter interface {
    write(b []byte) (*conn, os.Error)
    read(c *conn) *Reply
}

type Client struct {
    Addr     string
    Db       int
    Password string
    pool     *pool
}

// writes a command a and returns single the Reply object
func Send(rw ReaderWriter, args ...[]byte) *Reply {
    c, err := rw.write(buildCmd(args...))

    if err != nil {
        return &Reply{Err: err}
    }

    return rw.read(c)
}

// writes a command without calling read
func appendSend(rw ReaderWriter, args ...[]byte) (*conn, os.Error) {
    return rw.write(buildCmd(args...))
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

    return Send(rw, buf...)
}

// writes a command without calling read
func appendSendStr(rw ReaderWriter, name string, args ...string) (*conn, os.Error) {
    buf := make([][]byte, len(args)+1)
    buf[0] = []byte(name)

    for i, arg := range args {
        buf[i+1] = []byte(arg)
    }

    return appendSend(rw, buf...)
}

// creates a bytestring of the name and args parameters, then calls Send()
func SendStr(rw ReaderWriter, name string, args ...string) *Reply {
    c, err := appendSendStr(rw, name, args...)

    if err != nil {
        return &Reply{Err: err}
    }

    return rw.read(c)
}

// takes a [][]byte and returns a redis command formatted using
// the unified request protocol
func buildCmd(args ...[]byte) []byte {
    buf := bytes.NewBuffer(nil)

    buf.WriteByte(star)
    buf.WriteString(strconv.Itoa(len(args)))
    buf.Write(delim)

    for _, arg := range args {
        buf.WriteByte(dollar)
        buf.WriteString(strconv.Itoa(len(arg)))
        buf.Write(delim)
        buf.Write(arg)
        buf.Write(delim)
    }

    if logCmd {
        log.Printf("GODIS: %q", string(buf.Bytes()))
    }

    return buf.Bytes()
}

func New(addr string, db int, password string) *Client {
    if addr == "" {
        addr = "127.0.0.1:6379"
    }

    return &Client{Addr: addr, Db: db, Password: password, pool: newPool()}
}

func (c *Client) getConn() (*conn, os.Error) {
    cc := c.pool.pop()

    if cc != nil {
        return cc, nil
    }

    addr, err := net.ResolveTCPAddr(c.Addr)

    if err != nil {
        return nil, os.NewError("ResolveAddr error for " + c.Addr)
    }

    tcpc, err := net.DialTCP("tcp", nil, addr)
    if err != nil {
        return nil, os.NewError("Connection error " + addr.String())
    }

    cc = newConn(tcpc)
    err = cc.configConn(c)
    connCount++
    return cc, err
}

func (c *Client) read(conn *conn) *Reply {
    reply := conn.readReply()
    c.pool.push(conn)
    return reply
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

    if reply.Err != nil {
        p.end()
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
        p.end()
        return nil, err
    }

    p.appendMode = true
    p.replyCount++
    return p.conn, nil
}

func (p *Pipe) end() {
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
